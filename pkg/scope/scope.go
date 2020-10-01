package scope

import (
	"net/http"
	"regexp"
	"sync"
)

type Scope struct {
	mu    sync.Mutex
	rules []Rule
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

func New(rules []Rule) *Scope {
	s := &Scope{}
	if rules != nil {
		s.rules = rules
	}
	return s
}

func (s *Scope) Rules() []Rule {
	return s.rules
}

func (s *Scope) SetRules(rules []Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rules = rules
}

func (s *Scope) Match(req *http.Request, body []byte) bool {
	// TODO(?): Do we need to lock here as well?
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
