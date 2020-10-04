package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"

	// Register sqlite3 for use via database/sql.
	_ "github.com/mattn/go-sqlite3"
)

// Client implements reqlog.Repository.
type Client struct {
	db *sql.DB
}

// New returns a new Client.
func New(filename string) (*Client, error) {
	// Create directory for DB if it doesn't exist yet.
	if dbDir, _ := filepath.Split(filename); dbDir != "" {
		if _, err := os.Stat(dbDir); os.IsNotExist(err) {
			os.Mkdir(dbDir, 0755)
		}
	}

	opts := make(url.Values)
	opts.Set("_foreign_keys", "1")

	dsn := fmt.Sprintf("file:%v?%v", filename, opts.Encode())
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite: could not ping database: %v", err)
	}

	c := &Client{db: db}

	if err := c.prepareSchema(); err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare schema: %v", err)
	}

	return &Client{db: db}, nil
}

func (c Client) prepareSchema() error {
	_, err := c.db.Exec(`CREATE TABLE IF NOT EXISTS http_requests (
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

	_, err = c.db.Exec(`CREATE TABLE IF NOT EXISTS http_responses (
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

	_, err = c.db.Exec(`CREATE TABLE IF NOT EXISTS http_headers (
		id INTEGER PRIMARY KEY,
		req_id INTEGER REFERENCES http_requests(id) ON DELETE CASCADE,
		res_id INTEGER REFERENCES http_responses(id) ON DELETE CASCADE,
		key TEXT,
		value TEXT
	)`)
	if err != nil {
		return fmt.Errorf("could not create http_headers table: %v", err)
	}

	return nil
}

// Close uses the underlying database.
func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) FindRequestLogs(
	ctx context.Context,
	opts reqlog.FindRequestsOptions,
	scope *scope.Scope,
) (reqLogs []reqlog.Request, err error) {
	// TODO: Pass GraphQL field collections upstream, so we can query only
	// requested fields.
	// TODO: Use opts and scope to filter.
	reqQuery := `SELECT
		req.id,
		req.proto,
		req.url,
		req.method,
		req.body,
		req.timestamp,
		res.id,
		res.proto,
		res.status_code,
		res.status_reason,
		res.body,
		res.timestamp
	FROM http_requests req
	LEFT JOIN http_responses res ON req.id = res.req_id
	ORDER BY req.id DESC`

	rows, err := c.db.QueryContext(ctx, reqQuery)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reqLog reqlog.Request
		var resDTO httpResponse
		var statusReason *string
		var rawURL string

		err := rows.Scan(
			&reqLog.ID,
			&reqLog.Request.Proto,
			&rawURL,
			&reqLog.Request.Method,
			&reqLog.Body,
			&reqLog.Timestamp,
			&resDTO.ID,
			&resDTO.Proto,
			&resDTO.StatusCode,
			&statusReason,
			&resDTO.Body,
			&resDTO.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("sqlite: could not scan row: %v", err)
		}

		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, fmt.Errorf("sqlite: could not parse URL: %v", err)
		}
		reqLog.Request.URL = u

		if resDTO.ID != nil {
			status := strconv.Itoa(*resDTO.StatusCode) + " " + *statusReason
			reqLog.Response = &reqlog.Response{
				ID:        *resDTO.ID,
				RequestID: reqLog.ID,
				Response: http.Response{
					Status:     status,
					StatusCode: *resDTO.StatusCode,
					Proto:      *resDTO.Proto,
				},
				Body:      *resDTO.Body,
				Timestamp: *resDTO.Timestamp,
			}
		}

		reqLogs = append(reqLogs, reqLog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite: could not iterate over rows: %v", err)
	}
	rows.Close()

	reqHeadersStmt, err := c.db.PrepareContext(ctx, `SELECT key, value FROM http_headers WHERE req_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer reqHeadersStmt.Close()
	resHeadersStmt, err := c.db.PrepareContext(ctx, `SELECT key, value FROM http_headers WHERE res_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer resHeadersStmt.Close()

	for _, reqLog := range reqLogs {
		headers, err := findHeaders(ctx, reqHeadersStmt, reqLog.ID)
		if err != nil {
			return nil, fmt.Errorf("sqlite: could not query request headers: %v", err)
		}
		reqLog.Request.Header = headers

		if reqLog.Response != nil {
			headers, err := findHeaders(ctx, resHeadersStmt, reqLog.Response.ID)
			if err != nil {
				return nil, fmt.Errorf("sqlite: could not query response headers: %v", err)
			}
			reqLog.Response.Response.Header = headers
		}
	}

	return reqLogs, nil
}

func (c *Client) FindRequestLogByID(ctx context.Context, id int64) (reqlog.Request, error) {
	// TODO: Pass GraphQL field collections upstream, so we can query only
	// requested fields.
	reqQuery := `SELECT
		req.id,
		req.proto,
		req.url,
		req.method,
		req.body,
		req.timestamp,
		res.id,
		res.proto,
		res.status_code,
		res.status_reason,
		res.body,
		res.timestamp
	FROM http_requests req
	LEFT JOIN http_responses res ON req.id = res.req_id
	WHERE req_id = ?
	ORDER BY req.id DESC`

	var reqLog reqlog.Request
	var resDTO httpResponse
	var statusReason *string
	var rawURL string

	err := c.db.QueryRowContext(ctx, reqQuery, id).Scan(
		&reqLog.ID,
		&reqLog.Request.Proto,
		&rawURL,
		&reqLog.Request.Method,
		&reqLog.Body,
		&reqLog.Timestamp,
		&resDTO.ID,
		&resDTO.Proto,
		&resDTO.StatusCode,
		&statusReason,
		&resDTO.Body,
		&resDTO.Timestamp,
	)
	if err == sql.ErrNoRows {
		return reqlog.Request{}, reqlog.ErrRequestNotFound
	}
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not scan row: %v", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not parse URL: %v", err)
	}
	reqLog.Request.URL = u

	if resDTO.ID != nil {
		status := strconv.Itoa(*resDTO.StatusCode) + " " + *statusReason
		reqLog.Response = &reqlog.Response{
			ID:        *resDTO.ID,
			RequestID: reqLog.ID,
			Response: http.Response{
				Status:     status,
				StatusCode: *resDTO.StatusCode,
				Proto:      *resDTO.Proto,
			},
			Body:      *resDTO.Body,
			Timestamp: *resDTO.Timestamp,
		}
	}

	reqHeadersStmt, err := c.db.PrepareContext(ctx, `SELECT key, value FROM http_headers WHERE req_id = ?`)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer reqHeadersStmt.Close()
	resHeadersStmt, err := c.db.PrepareContext(ctx, `SELECT key, value FROM http_headers WHERE res_id = ?`)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not prepare statement: %v", err)
	}
	defer resHeadersStmt.Close()

	headers, err := findHeaders(ctx, reqHeadersStmt, reqLog.ID)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("sqlite: could not query request headers: %v", err)
	}
	reqLog.Request.Header = headers

	if reqLog.Response != nil {
		headers, err := findHeaders(ctx, resHeadersStmt, reqLog.Response.ID)
		if err != nil {
			return reqlog.Request{}, fmt.Errorf("sqlite: could not query response headers: %v", err)
		}
		reqLog.Response.Response.Header = headers
	}

	return reqLog, nil
}

func (c *Client) AddRequestLog(
	ctx context.Context,
	req http.Request,
	body []byte,
	timestamp time.Time,
) (*reqlog.Request, error) {

	reqLog := &reqlog.Request{
		Request:   req,
		Body:      body,
		Timestamp: timestamp,
	}

	tx, err := c.db.BeginTx(ctx, nil)
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
