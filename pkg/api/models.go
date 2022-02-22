package api

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/oklog/ulid"
)

func MarshalULID(u ulid.ULID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(u.String()))
	})
}

func UnmarshalULID(v interface{}) (ulid.ULID, error) {
	rawULID, ok := v.(string)
	if !ok {
		return ulid.ULID{}, fmt.Errorf("ulid must be a string")
	}

	u, err := ulid.Parse(rawULID)
	if err != nil {
		return ulid.ULID{}, fmt.Errorf("failed to parse ULID: %w", err)
	}

	return u, nil
}

func MarshalURL(u *url.URL) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(u.String()))
	})
}

func UnmarshalURL(v interface{}) (*url.URL, error) {
	rawURL, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("url must be a string")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	return u, nil
}
