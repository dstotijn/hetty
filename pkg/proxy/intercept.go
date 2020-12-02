package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	URL  string `yaml:"url"`
	Body string `yaml:"body"`
}

// TODO: create instance with responses as state

func NewInterceptor() error {
	var entries struct {
		Responses []Entry
	}

	path, _ := filepath.Abs("./pkg/proxy/intercept.yaml")
	fmt.Println(path)
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("intercept: error reading intercept.yaml: %v", err)
	}
	err = yaml.Unmarshal(fileContent, &entries)
	if err != nil {
		return fmt.Errorf("intercept: error parsing intercept.yaml: %v", err)
	}

	fmt.Println(entries.Responses)

	return nil
}

func ResponseInterceptor(next ResponseModifyFunc) ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		if res.Request.URL.String() == "http://www.ue.wroc.pl/" {
			//res.Header["X-Proxy"] = []string{"Hello"}

			fmt.Println("Intercepting ue.wroc.pl request!")

			err := changeBody(res, func(b []byte) []byte {
				//fmt.Println(string(b))

				prefix := []byte("<p>Pozdro poćwicz</p>")

				return append(prefix, b...)
			})
			if err != nil {
				panic(err)
			}

			//res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		return nil
	}
}
