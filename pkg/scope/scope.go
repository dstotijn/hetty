package scope

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"regexp"
	"sync"
)

type Scope struct {
	rules []Rule
	mu    sync.RWMutex
}

type Rule struct {
	URL    *regexp.Regexp
	Header Header
	Body   *regexp.Regexp
}

type Header struct {
	Key   *regexp.Regexp
	Value *regexp.Regexp
}

func (s *Scope) Rules() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rules
}

func (s *Scope) SetRules(rules []Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rules = rules
}

func (s *Scope) Match(req *http.Request, body []byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, rule := range s.rules {
		if matches := rule.Match(req, body); matches {
			return true
		}
	}

	return false
}

func (r Rule) Match(req *http.Request, body []byte) bool {
	if r.URL != nil {
		if matches := r.URL.MatchString(req.URL.String()); matches {
			return true
		}
	}

	for key, values := range req.Header {
		var keyMatches, valueMatches bool

		if r.Header.Key != nil {
			if matches := r.Header.Key.MatchString(key); matches {
				keyMatches = true
			}
		}

		if r.Header.Value != nil {
			for _, value := range values {
				if matches := r.Header.Value.MatchString(value); matches {
					valueMatches = true
					break
				}
			}
		}
		// When only key or value is set, match on whatever is set.
		// When both are set, both must match.
		switch {
		case r.Header.Key != nil && r.Header.Value == nil && keyMatches:
			return true
		case r.Header.Key == nil && r.Header.Value != nil && valueMatches:
			return true
		case r.Header.Key != nil && r.Header.Value != nil && keyMatches && valueMatches:
			return true
		}
	}

	if r.Body != nil {
		if matches := r.Body.Match(body); matches {
			return true
		}
	}

	return false
}

func regexpToString(r *regexp.Regexp) string {
	if r == nil {
		return ""
	}

	return r.String()
}

func stringToRegexp(s string) (*regexp.Regexp, error) {
	if s == "" {
		return nil, nil
	}

	return regexp.Compile(s)
}

type ruleDTO struct {
	URL    string
	Header struct {
		Key   string
		Value string
	}
	Body string
}

func (r Rule) MarshalBinary() ([]byte, error) {
	dto := ruleDTO{
		URL:  regexpToString(r.URL),
		Body: regexpToString(r.Body),
	}
	dto.Header.Key = regexpToString(r.Header.Key)
	dto.Header.Value = regexpToString(r.Header.Value)

	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(dto)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Rule) UnmarshalBinary(data []byte) error {
	dto := ruleDTO{}

	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&dto)
	if err != nil {
		return err
	}

	url, err := stringToRegexp(dto.URL)
	if err != nil {
		return err
	}

	headerKey, err := stringToRegexp(dto.Header.Key)
	if err != nil {
		return err
	}

	headerValue, err := stringToRegexp(dto.Header.Value)
	if err != nil {
		return err
	}

	body, err := stringToRegexp(dto.Body)
	if err != nil {
		return err
	}

	*r = Rule{
		URL: url,
		Header: Header{
			Key:   headerKey,
			Value: headerValue,
		},
		Body: body,
	}

	return nil
}
