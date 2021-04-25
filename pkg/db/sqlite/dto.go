package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dstotijn/hetty/pkg/reqlog"
)

type reqURL url.URL

type httpRequest struct {
	ID        int64     `db:"req_id"`
	Proto     string    `db:"req_proto"`
	URL       reqURL    `db:"url"`
	Method    string    `db:"method"`
	Body      []byte    `db:"req_body"`
	Timestamp time.Time `db:"req_timestamp"`
	httpResponse
}

type httpResponse struct {
	ID           sql.NullInt64  `db:"res_id"`
	RequestID    sql.NullInt64  `db:"res_req_id"`
	Proto        sql.NullString `db:"res_proto"`
	StatusCode   sql.NullInt64  `db:"status_code"`
	StatusReason sql.NullString `db:"status_reason"`
	Body         []byte         `db:"res_body"`
	Timestamp    sql.NullTime   `db:"res_timestamp"`
}

// Value implements driver.Valuer.
func (u *reqURL) Scan(value interface{}) error {
	rawURL, ok := value.(string)
	if !ok {
		return errors.New("sqlite: cannot scan non-string value")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("sqlite: could not parse URL: %w", err)
	}

	*u = reqURL(*parsed)

	return nil
}

func (dto httpRequest) toRequestLog() reqlog.Request {
	u := url.URL(dto.URL)
	reqLog := reqlog.Request{
		ID: dto.ID,
		Request: http.Request{
			Proto:  dto.Proto,
			Method: dto.Method,
			URL:    &u,
		},
		Body:      dto.Body,
		Timestamp: dto.Timestamp,
	}

	if dto.httpResponse.ID.Valid {
		reqLog.Response = &reqlog.Response{
			ID:        dto.httpResponse.ID.Int64,
			RequestID: dto.httpResponse.RequestID.Int64,
			Response: http.Response{
				Status:     strconv.FormatInt(dto.StatusCode.Int64, 10) + " " + dto.StatusReason.String,
				StatusCode: int(dto.StatusCode.Int64),
				Proto:      dto.httpResponse.Proto.String,
			},
			Body:      dto.httpResponse.Body,
			Timestamp: dto.httpResponse.Timestamp.Time,
		}
	}

	return reqLog
}
