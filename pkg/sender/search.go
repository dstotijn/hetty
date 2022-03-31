package sender

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/search"
)

var senderReqSearchKeyFns = map[string]func(req Request) string{
	"req.id":    func(req Request) string { return req.ID.String() },
	"req.proto": func(req Request) string { return req.Proto },
	"req.url": func(req Request) string {
		if req.URL == nil {
			return ""
		}
		return req.URL.String()
	},
	"req.method":    func(req Request) string { return req.Method },
	"req.body":      func(req Request) string { return string(req.Body) },
	"req.timestamp": func(req Request) string { return ulid.Time(req.ID.Time()).String() },
}

// TODO: Request and response headers search key functions.

// Matches returns true if the supplied search expression evaluates to true.
func (req Request) Matches(expr search.Expression) (bool, error) {
	switch e := expr.(type) {
	case search.PrefixExpression:
		return req.matchPrefixExpr(e)
	case search.InfixExpression:
		return req.matchInfixExpr(e)
	case search.StringLiteral:
		return req.matchStringLiteral(e)
	default:
		return false, fmt.Errorf("expression type (%T) not supported", expr)
	}
}

func (req Request) matchPrefixExpr(expr search.PrefixExpression) (bool, error) {
	switch expr.Operator {
	case search.TokOpNot:
		match, err := req.Matches(expr.Right)
		if err != nil {
			return false, err
		}

		return !match, nil
	default:
		return false, errors.New("operator is not supported")
	}
}

func (req Request) matchInfixExpr(expr search.InfixExpression) (bool, error) {
	switch expr.Operator {
	case search.TokOpAnd:
		left, err := req.Matches(expr.Left)
		if err != nil {
			return false, err
		}

		right, err := req.Matches(expr.Right)
		if err != nil {
			return false, err
		}

		return left && right, nil
	case search.TokOpOr:
		left, err := req.Matches(expr.Left)
		if err != nil {
			return false, err
		}

		right, err := req.Matches(expr.Right)
		if err != nil {
			return false, err
		}

		return left || right, nil
	}

	left, ok := expr.Left.(search.StringLiteral)
	if !ok {
		return false, errors.New("left operand must be a string literal")
	}

	leftVal := req.getMappedStringLiteral(left.Value)

	if leftVal == "req.headers" {
		match, err := search.MatchHTTPHeaders(expr.Operator, expr.Right, req.Header)
		if err != nil {
			return false, fmt.Errorf("failed to match request HTTP headers: %w", err)
		}

		return match, nil
	}

	if leftVal == "res.headers" && req.Response != nil {
		match, err := search.MatchHTTPHeaders(expr.Operator, expr.Right, req.Response.Header)
		if err != nil {
			return false, fmt.Errorf("failed to match response HTTP headers: %w", err)
		}

		return match, nil
	}

	if expr.Operator == search.TokOpRe || expr.Operator == search.TokOpNotRe {
		right, ok := expr.Right.(search.RegexpLiteral)
		if !ok {
			return false, errors.New("right operand must be a regular expression")
		}

		switch expr.Operator {
		case search.TokOpRe:
			return right.MatchString(leftVal), nil
		case search.TokOpNotRe:
			return !right.MatchString(leftVal), nil
		}
	}

	right, ok := expr.Right.(search.StringLiteral)
	if !ok {
		return false, errors.New("right operand must be a string literal")
	}

	rightVal := req.getMappedStringLiteral(right.Value)

	switch expr.Operator {
	case search.TokOpEq:
		return leftVal == rightVal, nil
	case search.TokOpNotEq:
		return leftVal != rightVal, nil
	case search.TokOpGt:
		// TODO(?) attempt to parse as int.
		return leftVal > rightVal, nil
	case search.TokOpLt:
		// TODO(?) attempt to parse as int.
		return leftVal < rightVal, nil
	case search.TokOpGtEq:
		// TODO(?) attempt to parse as int.
		return leftVal >= rightVal, nil
	case search.TokOpLtEq:
		// TODO(?) attempt to parse as int.
		return leftVal <= rightVal, nil
	default:
		return false, errors.New("unsupported operator")
	}
}

func (req Request) getMappedStringLiteral(s string) string {
	switch {
	case strings.HasPrefix(s, "req."):
		fn, ok := senderReqSearchKeyFns[s]
		if ok {
			return fn(req)
		}
	case strings.HasPrefix(s, "res."):
		if req.Response == nil {
			return ""
		}

		fn, ok := reqlog.ResLogSearchKeyFns[s]
		if ok {
			return fn(*req.Response)
		}
	}

	return s
}

func (req Request) matchStringLiteral(strLiteral search.StringLiteral) (bool, error) {
	for _, fn := range senderReqSearchKeyFns {
		if strings.Contains(
			strings.ToLower(fn(req)),
			strings.ToLower(strLiteral.Value),
		) {
			return true, nil
		}
	}

	if req.Response != nil {
		for _, fn := range reqlog.ResLogSearchKeyFns {
			if strings.Contains(
				strings.ToLower(fn(*req.Response)),
				strings.ToLower(strLiteral.Value),
			) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (req Request) MatchScope(s *scope.Scope) bool {
	for _, rule := range s.Rules() {
		if rule.URL != nil && req.URL != nil {
			if matches := rule.URL.MatchString(req.URL.String()); matches {
				return true
			}
		}

		for key, values := range req.Header {
			var keyMatches, valueMatches bool

			if rule.Header.Key != nil {
				if matches := rule.Header.Key.MatchString(key); matches {
					keyMatches = true
				}
			}

			if rule.Header.Value != nil {
				for _, value := range values {
					if matches := rule.Header.Value.MatchString(value); matches {
						valueMatches = true
						break
					}
				}
			}
			// When only key or value is set, match on whatever is set.
			// When both are set, both must match.
			switch {
			case rule.Header.Key != nil && rule.Header.Value == nil && keyMatches:
				return true
			case rule.Header.Key == nil && rule.Header.Value != nil && valueMatches:
				return true
			case rule.Header.Key != nil && rule.Header.Value != nil && keyMatches && valueMatches:
				return true
			}
		}

		if rule.Body != nil {
			if matches := rule.Body.Match(req.Body); matches {
				return true
			}
		}
	}

	return false
}
