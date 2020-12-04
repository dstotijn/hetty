package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func changeBody(res *http.Response, modifer func(body []byte) []byte) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Could not read response body: %v", err)
	}

	contentEncoding := res.Header.Get("Content-Encoding")

	if contentEncoding == "" {
		newBody := modifer(body)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
	}

	if contentEncoding == "gzip" {
		// TMP!
		//res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("Could not create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		body, err = ioutil.ReadAll(gzipReader)

		// TODO: Gzip this body
		newBody := modifer(body)
		res.Header.Set("Content-Encoding", "")
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))

		if err != nil {
			return fmt.Errorf("Could not read gzipped response body: %v", err)
		}
	}

	return nil
}

type Entry struct {
	UrlEquals     string `yaml:"url_equals"`
	UrlStartsWith string `yaml:"url_starts_with"`
	UrlEndsWith   string `yaml:"url_ends_with"`
	Body          string
}

// TODO: create instance with responses as state

type Intercept struct {
	responses []Entry
}

func NewIntercep() (*Intercept, error) {
	var entries struct {
		Responses []Entry
	}

	path, _ := filepath.Abs("./pkg/proxy/intercept.yaml")
	fmt.Println(path)
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("intercept: error reading intercept.yaml: %v", err)
	}
	err = yaml.Unmarshal(fileContent, &entries)
	if err != nil {
		return nil, fmt.Errorf("intercept: error parsing intercept.yaml: %v", err)
	}

	return &Intercept{responses: entries.Responses}, nil
}

func (intercept *Intercept) ResponseInterceptorFromYAML(next ResponseModifyFunc) ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		for _, response := range intercept.responses {
			replaceBody := func() {
				err := changeBody(res, func(b []byte) []byte {
					return []byte(response.Body)
				})
				if err != nil {
					panic(err)
				}
			}

			url := res.Request.URL.String()

			if url == response.UrlEquals {
				replaceBody()
			}

			if response.UrlStartsWith != "" && strings.HasPrefix(url, response.UrlStartsWith) {
				replaceBody()
			}

			if response.UrlEndsWith != "" && strings.HasSuffix(url, response.UrlEndsWith) {
				replaceBody()
			}
		}

		return nil
	}
}
