package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/dstotijn/hetty/pkg/proj"
)

const moduleName = "scope"

type Scope struct {
	rules []Rule
	repo  Repository

	mu sync.RWMutex
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

func New(repo Repository, projService *proj.Service) *Scope {
	s := &Scope{
		repo: repo,
	}

	projService.OnProjectOpen(func(_ string) error {
		err := s.load(context.Background())
		if err == proj.ErrNoSettings {
			return nil
		}
		if err != nil {
			return fmt.Errorf("scope: could not load scope: %v", err)
		}
		return nil
	})
	projService.OnProjectClose(func(_ string) error {
		s.unload()
		return nil
	})

	return s
}

func (s *Scope) Rules() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rules
}

func (s *Scope) load(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var rules []Rule
	err := s.repo.FindSettingsByModule(ctx, moduleName, &rules)
	if err == proj.ErrNoSettings {
		return err
	}
	if err != nil {
		return fmt.Errorf("scope: could not load scope settings: %v", err)
	}

	s.rules = rules

	return nil
}

func (s *Scope) unload() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = nil
}

func (s *Scope) SetRules(ctx context.Context, rules []Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.UpsertSettings(ctx, moduleName, rules); err != nil {
		return fmt.Errorf("scope: cannot set rules in repository: %v", err)
	}

	s.rules = rules

	return nil
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

// MarshalJSON implements json.Marshaler.
func (r Rule) MarshalJSON() ([]byte, error) {
	type headerDTO struct {
		Key   string
		Value string
	}
	type ruleDTO struct {
		URL    string
		Header headerDTO
		Body   string
	}

	dto := ruleDTO{
		URL: regexpToString(r.URL),
		Header: headerDTO{
			Key:   regexpToString(r.Header.Key),
			Value: regexpToString(r.Header.Value),
		},
		Body: regexpToString(r.Body),
	}

	return json.Marshal(dto)
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Rule) UnmarshalJSON(data []byte) error {
	type headerDTO struct {
		Key   string
		Value string
	}
	type ruleDTO struct {
		URL    string
		Header headerDTO
		Body   string
	}

	var dto ruleDTO
	if err := json.Unmarshal(data, &dto); err != nil {
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
