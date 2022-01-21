package api

import (
	"fmt"
	"io"
	"strconv"

	"github.com/oklog/ulid"
)

type ULID ulid.ULID

func (u *ULID) UnmarshalGQL(v interface{}) (err error) {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("ulid must be a string")
	}

	id, err := ulid.Parse(str)
	if err != nil {
		return fmt.Errorf("failed to parse ULID: %w", err)
	}

	*u = ULID(id)

	return nil
}

func (u ULID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(ulid.ULID(u).String()))
}
