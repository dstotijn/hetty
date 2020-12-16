package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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
	} else if contentEncoding == "gzip" {
		// TMP!
		//res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("Could not create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		body, err = ioutil.ReadAll(gzipReader)

		if err != nil {
			return fmt.Errorf("Could not read gzipped response body: %v", err)
		}

		newBody := modifer(body)
		var gzipBodyBuffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzipBodyBuffer)
		if _, err := gzipWriter.Write(newBody); err != nil {
			return fmt.Errorf("Could not write gzip body: %v", err)
		}
		if err := gzipWriter.Close(); err != nil {
			return fmt.Errorf("Could not close gzip body writer: %v", err)
		}

		res.Body = ioutil.NopCloser(&gzipBodyBuffer)
		res.Header.Set("Content-Length", strconv.Itoa(gzipBodyBuffer.Len()))

	} else if contentEncoding == "br" {
		// TODO: Handle brotli properly
		newBody := modifer(body)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
		res.Header.Del("Content-Encoding")
	} else {
		return fmt.Errorf("Unknown encoding: %v", contentEncoding)
	}

	return nil
}

type RequestEntry struct {
	UrlEquals     string `yaml:"url_equals"`
	UrlStartsWith string `yaml:"url_starts_with"`
	UrlEndsWith   string `yaml:"url_ends_with"`
	Body          string
	Headers       map[string]string `yaml:"headers,omitempty"`
}

type ResponseEntry struct {
	UrlEquals     string `yaml:"url_equals"`
	UrlStartsWith string `yaml:"url_starts_with"`
	UrlEndsWith   string `yaml:"url_ends_with"`
	Body          string
	Status        int
	Headers       map[string]string `yaml:"headers,omitempty"`
}

type Intercept struct {
	getRequests  func() ([]RequestEntry, error)
	getResponses func() ([]ResponseEntry, error)
}

func NewIntercept(
	getRequests func() ([]RequestEntry, error),
	getResponses func() ([]ResponseEntry, error),
) (*Intercept, error) {
	return &Intercept{
		getRequests:  getRequests,
		getResponses: getResponses,
	}, nil
}

func GetRequestsFromYaml() ([]RequestEntry, error) {
	return []RequestEntry{}, nil
}

func GetResponsesFromYaml() ([]ResponseEntry, error) {
	var responses []ResponseEntry
	path, _ := filepath.Abs("./intercept.yaml")
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("intercept: error reading intercept.yaml: %v", err)
	}
	err = yaml.Unmarshal(fileContent, struct{ Responses *[]ResponseEntry }{Responses: &responses})
	if err != nil {
		return nil, fmt.Errorf("intercept: error parsing intercept.yaml: %v", err)
	}
	return responses, nil
}

func (intercept *Intercept) RequestInterceptor(next RequestModifyFunc) RequestModifyFunc {
	return func(req *http.Request) {
		next(req)

		//if err := next(req); err != nil {
		//	log.Fatal(err)
		//}

		requests, err := intercept.getRequests()
		if err != nil {
			log.Fatal(err)
		}

		for _, request := range requests {

			url := req.URL.String()

			if url == request.UrlEquals ||
				(request.UrlStartsWith != "" && strings.HasPrefix(url, request.UrlStartsWith)) ||
				(request.UrlEndsWith != "" && strings.HasSuffix(url, request.UrlEndsWith)) {

				//if request.Body != "" {
				//	err := changeBody(req, func(b []byte) []byte {
				//		return []byte(request.Body)
				//	})
				//	if err != nil {
				//		panic(err)
				//	}
				//}

				if len(request.Headers) != 0 {
					for key, value := range request.Headers {
						req.Header.Set(key, value)
					}
				}

				//if request.Status != 0 {
				//  res.StatusCode = request.Status
				//}
			}
		}
	}
}

func (intercept *Intercept) ResponseInterceptor(next ResponseModifyFunc) ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		responses, err := intercept.getResponses()
		if err != nil {
			return err
		}

		for _, response := range responses {

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
