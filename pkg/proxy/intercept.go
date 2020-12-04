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

type ResponseEntry struct {
	UrlEquals     string `yaml:"url_equals"`
	UrlStartsWith string `yaml:"url_starts_with"`
	UrlEndsWith   string `yaml:"url_ends_with"`
	Body          string
	Status        int
	Headers       map[string]string `yaml:"headers,omitempty"`
}

// TODO: create instance with responses as state

type Intercept struct {
	responses []ResponseEntry
}

func NewIntercep() (*Intercept, error) {
	return &Intercept{}, nil
}

func (intercept *Intercept) updateYamlEntries() error {
	path, _ := filepath.Abs("./pkg/proxy/intercept.yaml")
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("intercept: error reading intercept.yaml: %v", err)
	}
	err = yaml.Unmarshal(fileContent, struct{ Responses *[]ResponseEntry }{Responses: &intercept.responses})
	if err != nil {
		return fmt.Errorf("intercept: error parsing intercept.yaml: %v", err)
	}
	return nil
}

func (intercept *Intercept) ResponseInterceptorFromYaml(next ResponseModifyFunc) ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		err := intercept.updateYamlEntries()
		if err != nil {
			return err
		}

		for _, response := range intercept.responses {

			url := res.Request.URL.String()

			if url == response.UrlEquals ||
				(response.UrlStartsWith != "" && strings.HasPrefix(url, response.UrlStartsWith)) ||
				(response.UrlEndsWith != "" && strings.HasSuffix(url, response.UrlEndsWith)) {

				if response.Body != "" {
					err := changeBody(res, func(b []byte) []byte {
						return []byte(response.Body)
					})
					if err != nil {
						panic(err)
					}
				}

				if len(response.Headers) != 0 {
					for key, value := range response.Headers {
						res.Header.Set(key, value)
					}
				}

				if response.Status != 0 {
					res.StatusCode = response.Status
				}
			}
		}

		return nil
	}
}
