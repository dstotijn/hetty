package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"

	"github.com/99designs/gqlgen/graphql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

var regexpFn = func(pattern string, value interface{}) (bool, error) {
	switch v := value.(type) {
	case string:
		return regexp.MatchString(pattern, v)
	case int64:
		return regexp.MatchString(pattern, string(v))
	case []byte:
		return regexp.Match(pattern, v)
	default:
		return false, fmt.Errorf("unsupported type %T", v)
	}
}

// Client implements reqlog.Repository.
type Client struct {
	db            *sqlx.DB
	dbPath        string
	activeProject string
}

type httpRequestLogsQuery struct {
	requestCols        []string
	requestHeaderCols  []string
	responseHeaderCols []string
	joinResponse       bool
}

func init() {
	sql.Register("sqlite3_with_regexp", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("regexp", regexpFn, false); err != nil {
				return err
			}
			return nil
		},
	})
}

func New(dbPath string) (*Client, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dbPath, 0755); err != nil {
			return nil, fmt.Errorf("proj: could not create project directory: %v", err)
		}
	}
	return &Client{
		dbPath: dbPath,
	}, nil
}

// OpenProject opens a project database.
func (c *Client) OpenProject(name string) error {
	if c.db != nil {
		return errors.New("sqlite: there is already a project open")
	}

	opts := make(url.Values)
	opts.Set("_foreign_keys", "1")

	dbPath := filepath.Join(c.dbPath, name+".db")
	dsn := fmt.Sprintf("file:%v?%v", dbPath, opts.Encode())
	db, err := sqlx.Open("sqlite3_with_regexp", dsn)
	if err != nil {
		return fmt.Errorf("sqlite: could not open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("sqlite: could not ping database: %v", err)
	}

	if err := prepareSchema(db); err != nil {
		return fmt.Errorf("sqlite: could not prepare schema: %v", err)
	}

	c.db = db
	c.activeProject = name

	return nil
}

func (c *Client) Projects() ([]proj.Project, error) {
	files, err := ioutil.ReadDir(c.dbPath)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not read projects directory: %v", err)
	}

	projects := make([]proj.Project, len(files))
	for i, file := range files {
		projName := strings.TrimSuffix(file.Name(), ".db")
		projects[i] = proj.Project{
			Name:     projName,
			IsActive: c.activeProject == projName,
		}
	}

	return projects, nil
}

func prepareSchema(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS http_requests (
		id INTEGER PRIMARY KEY,
		proto TEXT,
		url TEXT,
		method TEXT,
		body BLOB,
		timestamp DATETIME
	)`)
	if err != nil {
		return fmt.Errorf("could not create http_requests table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS http_responses (
		id INTEGER PRIMARY KEY,
		req_id INTEGER REFERENCES http_requests(id) ON DELETE CASCADE,
		proto TEXT,
		status_code INTEGER,
		status_reason TEXT,
		body BLOB,
		timestamp DATETIME
	)`)
	if err != nil {
		return fmt.Errorf("could not create http_responses table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS http_headers (
		id INTEGER PRIMARY KEY,
		req_id INTEGER REFERENCES http_requests(id) ON DELETE CASCADE,
		res_id INTEGER REFERENCES http_responses(id) ON DELETE CASCADE,
		key TEXT,
		value TEXT
	)`)
	if err != nil {
		return fmt.Errorf("could not create http_headers table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS settings (
		module TEXT PRIMARY KEY,
		settings TEXT
	)`)
	if err != nil {
		return fmt.Errorf("could not create settings table: %v", err)
	}

	return nil
}

// Close uses the underlying database if it's open.
func (c *Client) Close() error {
	if c.db == nil {
		return nil
	}
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("sqlite: could not close database: %v", err)
	}

	c.db = nil
	c.activeProject = ""

	return nil
}

func (c *Client) DeleteProject(name string) error {
	if err := os.Remove(filepath.Join(c.dbPath, name+".db")); err != nil {
		return fmt.Errorf("sqlite: could not remove database file: %v", err)
	}

	return nil
}

var reqFieldToColumnMap = map[string]string{
	"proto":     "proto AS req_proto",
	"url":       "url",
	"method":    "method",
	"body":      "body AS req_body",
	"timestamp": "timestamp AS req_timestamp",
}

var resFieldToColumnMap = map[string]string{
	"requestId":    "req_id AS res_req_id",
	"proto":        "proto AS res_proto",
	"statusCode":   "status_code",
	"statusReason": "status_reason",
	"body":         "body AS res_body",
	"timestamp":    "timestamp AS res_timestamp",
}

var headerFieldToColumnMap = map[string]string{
	"key":   "key",
	"value": "value",
}

func (c *Client) ClearRequestLogs(ctx context.Context) error {
	if c.db == nil {
		return proj.ErrNoProject
	}
	_, err := c.db.Exec("DELETE FROM http_requests")
	if err != nil {
		return fmt.Errorf("sqlite: could not delete requests: %v", err)
	}
	return nil
}

func (c *Client) FindRequestLogs(
	ctx context.Context,
	filter reqlog.FindRequestsFilter,
	scope *scope.Scope,
) (reqLogs []reqlog.Request, err error) {
	if c.db == nil {
		return nil, proj.ErrNoProject
	}

	httpReqLogsQuery := parseHTTPRequestLogsQuery(ctx)

	reqQuery := sq.
		Select(httpReqLogsQuery.requestCols...).
		From("http_requests req").
		OrderBy("req.id DESC")
	if httpReqLogsQuery.joinResponse {
		reqQuery = reqQuery.LeftJoin("http_responses res ON req.id = res.req_id")
	}

	if filter.OnlyInScope && scope != nil {
		var ruleExpr []sq.Sqlizer
		for _, rule := range scope.Rules() {
			if rule.URL != nil {
				ruleExpr = append(ruleExpr, sq.Expr("regexp(?, req.url)", rule.URL.String()))
			}
		}
		if len(ruleExpr) > 0 {
			reqQuery = reqQuery.Where(sq.Or(ruleExpr))
		}
	}

	sql, args, err := reqQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not parse query: %v", err)
	}

	rows, err := c.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dto httpRequest
		err = rows.StructScan(&dto)
		if err != nil {
			return nil, fmt.Errorf("sqlite: could not scan row: %v", err)
		}
		reqLogs = append(reqLogs, dto.toRequestLog())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite: could not iterate over rows: %v", err)
	}
	rows.Close()

	if err := c.queryHeaders(ctx, httpReqLogsQuery, reqLogs); err != nil {
		return nil, fmt.Errorf("sqlite: could not query headers: %v", err)
	}

	return reqLogs, nil
}

func (c *Client) FindRequestLogByID(ctx context.Context, id int64) (reqlog.Request, error) {
	if c.db == nil {
		return reqlog.Request{}, proj.ErrNoProject
	}
	httpReqLogsQuery := parseHTTPRequestLogsQuery(ctx)

	reqQuery := sq.
		Select(httpReqLogsQuery.requestCols...).
		From("http_requests req").
		Where("req.id = ?")
	if httpReqLogsQuery.joinResponse {
		reqQuery = reqQuery.LeftJoin("http_responses res ON req.id = res.req_id")
	}

	reqSQL, _, err := reqQuery.ToSql()
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not parse query: %v", err)
	}

	row := c.db.QueryRowxContext(ctx, reqSQL, id)
	var dto httpRequest
	err = row.StructScan(&dto)
	if err == sql.ErrNoRows {
		return reqlog.Request{}, reqlog.ErrRequestNotFound
	}
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not scan row: %v", err)
	}
	reqLog := dto.toRequestLog()

	reqLogs := []reqlog.Request{reqLog}
	if err := c.queryHeaders(ctx, httpReqLogsQuery, reqLogs); err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not query headers: %v", err)
	}

	return reqLogs[0], nil
}

func (c *Client) AddRequestLog(
	ctx context.Context,
	req http.Request,
	body []byte,
	timestamp time.Time,
) (*reqlog.Request, error) {
	if c.db == nil {
		return nil, proj.ErrNoProject
	}

	reqLog := &reqlog.Request{
		Request:   req,
		Body:      body,
		Timestamp: timestamp,
	}

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not start transaction: %v", err)
	}
	defer tx.Rollback()

	reqStmt, err := tx.PrepareContext(ctx, `INSERT INTO http_requests (
		proto,
		url,
		method,
		body,
		timestamp
	) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer reqStmt.Close()

	result, err := reqStmt.ExecContext(ctx,
		reqLog.Request.Proto,
		reqLog.Request.URL.String(),
		reqLog.Request.Method,
		reqLog.Body,
		reqLog.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not execute statement: %v", err)
	}

	reqID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not get last insert ID: %v", err)
	}
	reqLog.ID = reqID

	headerStmt, err := tx.PrepareContext(ctx, `INSERT INTO http_headers (
		req_id,
		key,
		value
	) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer headerStmt.Close()

	err = insertHeaders(ctx, headerStmt, reqID, reqLog.Request.Header)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not insert http headers: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("sqlite: could not commit transaction: %v", err)
	}

	return reqLog, nil
}

func (c *Client) AddResponseLog(
	ctx context.Context,
	reqID int64,
	res http.Response,
	body []byte,
	timestamp time.Time,
) (*reqlog.Response, error) {
	if c.db == nil {
		return nil, proj.ErrNoProject
	}

	resLog := &reqlog.Response{
		RequestID: reqID,
		Response:  res,
		Body:      body,
		Timestamp: timestamp,
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not start transaction: %v", err)
	}
	defer tx.Rollback()

	resStmt, err := tx.PrepareContext(ctx, `INSERT INTO http_responses (
		req_id,
		proto,
		status_code,
		status_reason,
		body,
		timestamp
	) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer resStmt.Close()

	var statusReason string
	if len(resLog.Response.Status) > 4 {
		statusReason = resLog.Response.Status[4:]
	}

	result, err := resStmt.ExecContext(ctx,
		resLog.RequestID,
		resLog.Response.Proto,
		resLog.Response.StatusCode,
		statusReason,
		resLog.Body,
		resLog.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not execute statement: %v", err)
	}

	resID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not get last insert ID: %v", err)
	}
	resLog.ID = resID

	headerStmt, err := tx.PrepareContext(ctx, `INSERT INTO http_headers (
		res_id,
		key,
		value
	) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer headerStmt.Close()

	err = insertHeaders(ctx, headerStmt, resID, resLog.Response.Header)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not insert http headers: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("sqlite: could not commit transaction: %v", err)
	}

	return resLog, nil
}

func (c *Client) UpsertSettings(ctx context.Context, module string, settings interface{}) error {
	if c.db == nil {
		// TODO: Fix where `ErrNoProject` lives.
		return proj.ErrNoProject
	}

	jsonSettings, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("sqlite: could not encode settings as JSON: %v", err)
	}

	_, err = c.db.ExecContext(ctx,
		`INSERT INTO settings (module, settings) VALUES (?, ?)
		ON CONFLICT(module) DO UPDATE SET settings = ?`, module, jsonSettings, jsonSettings)
	if err != nil {
		return fmt.Errorf("sqlite: could not insert scope settings: %v", err)
	}

	return nil
}

func (c *Client) FindSettingsByModule(ctx context.Context, module string, settings interface{}) error {
	if c.db == nil {
		return proj.ErrNoProject
	}

	var jsonSettings []byte
	row := c.db.QueryRowContext(ctx, `SELECT settings FROM settings WHERE module = ?`, module)
	err := row.Scan(&jsonSettings)
	if err == sql.ErrNoRows {
		return proj.ErrNoSettings
	}
	if err != nil {
		return fmt.Errorf("sqlite: could not scan row: %v", err)
	}

	if err := json.Unmarshal(jsonSettings, &settings); err != nil {
		return fmt.Errorf("sqlite: could not decode settings from JSON: %v", err)
	}

	return nil
}

func insertHeaders(ctx context.Context, stmt *sql.Stmt, id int64, headers http.Header) error {
	for key, values := range headers {
		for _, value := range values {
			if _, err := stmt.ExecContext(ctx, id, key, value); err != nil {
				return fmt.Errorf("could not execute statement: %v", err)
			}
		}
	}
	return nil
}

func findHeaders(ctx context.Context, stmt *sql.Stmt, id int64) (http.Header, error) {
	headers := make(http.Header)
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		err := rows.Scan(
			&key,
			&value,
		)
		if err != nil {
			return nil, fmt.Errorf("sqlite: could not scan row: %v", err)
		}
		headers.Add(key, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite: could not iterate over rows: %v", err)
	}

	return headers, nil
}

func parseHTTPRequestLogsQuery(ctx context.Context) httpRequestLogsQuery {
	var joinResponse bool
	var reqHeaderCols, resHeaderCols []string

	opCtx := graphql.GetOperationContext(ctx)
	reqFields := graphql.CollectFieldsCtx(ctx, nil)
	reqCols := []string{"req.id AS req_id", "res.id AS res_id"}

	for _, reqField := range reqFields {
		if col, ok := reqFieldToColumnMap[reqField.Name]; ok {
			reqCols = append(reqCols, "req."+col)
		}
		if reqField.Name == "headers" {
			headerFields := graphql.CollectFields(opCtx, reqField.Selections, nil)
			for _, headerField := range headerFields {
				if col, ok := headerFieldToColumnMap[headerField.Name]; ok {
					reqHeaderCols = append(reqHeaderCols, col)
				}
			}
		}
		if reqField.Name == "response" {
			joinResponse = true
			resFields := graphql.CollectFields(opCtx, reqField.Selections, nil)
			for _, resField := range resFields {
				if resField.Name == "headers" {
					reqCols = append(reqCols, "res.id AS res_id")
					headerFields := graphql.CollectFields(opCtx, resField.Selections, nil)
					for _, headerField := range headerFields {
						if col, ok := headerFieldToColumnMap[headerField.Name]; ok {
							resHeaderCols = append(resHeaderCols, col)
						}
					}
				}
				if col, ok := resFieldToColumnMap[resField.Name]; ok {
					reqCols = append(reqCols, "res."+col)
				}
			}
		}
	}

	return httpRequestLogsQuery{
		requestCols:        reqCols,
		requestHeaderCols:  reqHeaderCols,
		responseHeaderCols: resHeaderCols,
		joinResponse:       joinResponse,
	}
}

func (c *Client) queryHeaders(
	ctx context.Context,
	query httpRequestLogsQuery,
	reqLogs []reqlog.Request,
) error {
	if len(query.requestHeaderCols) > 0 {
		reqHeadersQuery, _, err := sq.
			Select(query.requestHeaderCols...).
			From("http_headers").Where("req_id = ?").
			ToSql()
		if err != nil {
			return fmt.Errorf("could not parse request headers query: %v", err)
		}
		reqHeadersStmt, err := c.db.PrepareContext(ctx, reqHeadersQuery)
		if err != nil {
			return fmt.Errorf("could not prepare statement: %v", err)
		}
		defer reqHeadersStmt.Close()
		for i := range reqLogs {
			headers, err := findHeaders(ctx, reqHeadersStmt, reqLogs[i].ID)
			if err != nil {
				return fmt.Errorf("could not query request headers: %v", err)
			}
			reqLogs[i].Request.Header = headers
		}
	}

	if len(query.responseHeaderCols) > 0 {
		resHeadersQuery, _, err := sq.
			Select(query.responseHeaderCols...).
			From("http_headers").Where("res_id = ?").
			ToSql()
		if err != nil {
			return fmt.Errorf("could not parse response headers query: %v", err)
		}
		resHeadersStmt, err := c.db.PrepareContext(ctx, resHeadersQuery)
		if err != nil {
			return fmt.Errorf("could not prepare statement: %v", err)
		}
		defer resHeadersStmt.Close()
		for i := range reqLogs {
			if reqLogs[i].Response == nil {
				continue
			}
			headers, err := findHeaders(ctx, resHeadersStmt, reqLogs[i].Response.ID)
			if err != nil {
				return fmt.Errorf("could not query response headers: %v", err)
			}
			reqLogs[i].Response.Response.Header = headers
		}
	}

	return nil
}

func (c *Client) IsOpen() bool {
	return c.db != nil
}
