package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func gunzipResponseBody(res *http.Response) error {
	if res.Header.Get("Content-Encoding") != "gzip" {
		return nil
	}

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return fmt.Errorf("proxy: could not create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	buf := &bytes.Buffer{}

	//nolint:gosec
	if _, err := io.Copy(buf, gzipReader); err != nil {
		return fmt.Errorf("proxy: could not read gzipped response body: %w", err)
	}

	res.Body = io.NopCloser(buf)
	res.Header.Del("Content-Encoding")
	res.Header.Set("Content-Length", fmt.Sprint(buf.Len()))
	res.ContentLength = int64(buf.Len())

	return nil
}
