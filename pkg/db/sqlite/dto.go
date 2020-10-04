package sqlite

import "time"

type httpResponse struct {
	ID         *int64
	Proto      *string
	StatusCode *int
	Body       *[]byte
	Timestamp  *time.Time
}
