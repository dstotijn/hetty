package intercept

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dstotijn/hetty/pkg/filter"
	"github.com/dstotijn/hetty/pkg/scope"
)

//nolint:unparam
var reqFilterKeyFns = map[string]func(req *http.Request) (string, error){
	"proto": func(req *http.Request) (string, error) { return req.Proto, nil },
	"url": func(req *http.Request) (string, error) {
		if req.URL == nil {
			return "", nil
		}
		return req.URL.String(), nil
	},
	"method": func(req *http.Request) (string, error) { return req.Method, nil },
	"body": func(req *http.Request) (string, error) {
		if req.Body == nil {
			return "", nil
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			return "", err
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return string(body), nil
	},
}

//nolint:unparam
var resFilterKeyFns = map[string]func(res *http.Response) (string, error){
	"proto":      func(res *http.Response) (string, error) { return res.Proto, nil },
	"statusCode": func(res *http.Response) (string, error) { return strconv.Itoa(res.StatusCode), nil },
	"statusReason": func(res *http.Response) (string, error) {
		statusReasonSubs := strings.SplitN(res.Status, " ", 2)

		if len(statusReasonSubs) != 2 {
			return "", fmt.Errorf("invalid response status %q", res.Status)
		}
		return statusReasonSubs[1], nil
	},
	"body": func(res *http.Response) (string, error) {
		if res.Body == nil {
			return "", nil
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		return string(body), nil
	},
}

// MatchRequestFilter returns true if an HTTP request matches the request filter expression.
func MatchRequestFilter(req *http.Request, expr filter.Expression) (bool, error) {
	switch e := expr.(type) {
	case filter.PrefixExpression:
		return matchReqPrefixExpr(req, e)
	case filter.InfixExpression:
		return matchReqInfixExpr(req, e)
	case filter.StringLiteral:
		return matchReqStringLiteral(req, e)
	default:
		return false, fmt.Errorf("expression type (%T) not supported", expr)
	}
}

func matchReqPrefixExpr(req *http.Request, expr filter.PrefixExpression) (bool, error) {
	switch expr.Operator {
	case filter.TokOpNot:
		match, err := MatchRequestFilter(req, expr.Right)
		if err != nil {
			return false, err
		}

		return !match, nil
	default:
		return false, errors.New("operator is not supported")
	}
}

func matchReqInfixExpr(req *http.Request, expr filter.InfixExpression) (bool, error) {
	switch expr.Operator {
	case filter.TokOpAnd:
		left, err := MatchRequestFilter(req, expr.Left)
		if err != nil {
			return false, err
		}

		right, err := MatchRequestFilter(req, expr.Right)
		if err != nil {
			return false, err
		}

		return left && right, nil
	case filter.TokOpOr:
		left, err := MatchRequestFilter(req, expr.Left)
		if err != nil {
			return false, err
		}

		right, err := MatchRequestFilter(req, expr.Right)
		if err != nil {
			return false, err
		}

		return left || right, nil
	}

	left, ok := expr.Left.(filter.StringLiteral)
	if !ok {
		return false, errors.New("left operand must be a string literal")
	}

	leftVal, err := getMappedStringLiteralFromReq(req, left.Value)
	if err != nil {
		return false, fmt.Errorf("failed to get string literal from request for left operand: %w", err)
	}

	if leftVal == "headers" {
		match, err := filter.MatchHTTPHeaders(expr.Operator, expr.Right, req.Header)
		if err != nil {
			return false, fmt.Errorf("failed to match request HTTP headers: %w", err)
		}

		return match, nil
	}

	if expr.Operator == filter.TokOpRe || expr.Operator == filter.TokOpNotRe {
		right, ok := expr.Right.(filter.RegexpLiteral)
		if !ok {
			return false, errors.New("right operand must be a regular expression")
		}

		switch expr.Operator {
		case filter.TokOpRe:
			return right.MatchString(leftVal), nil
		case filter.TokOpNotRe:
			return !right.MatchString(leftVal), nil
		}
	}

	right, ok := expr.Right.(filter.StringLiteral)
	if !ok {
		return false, errors.New("right operand must be a string literal")
	}

	rightVal, err := getMappedStringLiteralFromReq(req, right.Value)
	if err != nil {
		return false, fmt.Errorf("failed to get string literal from request for right operand: %w", err)
	}

	switch expr.Operator {
	case filter.TokOpEq:
		return leftVal == rightVal, nil
	case filter.TokOpNotEq:
		return leftVal != rightVal, nil
	case filter.TokOpGt:
		// TODO(?) attempt to parse as int.
		return leftVal > rightVal, nil
	case filter.TokOpLt:
		// TODO(?) attempt to parse as int.
		return leftVal < rightVal, nil
	case filter.TokOpGtEq:
		// TODO(?) attempt to parse as int.
		return leftVal >= rightVal, nil
	case filter.TokOpLtEq:
		// TODO(?) attempt to parse as int.
		return leftVal <= rightVal, nil
	default:
		return false, errors.New("unsupported operator")
	}
}

func getMappedStringLiteralFromReq(req *http.Request, s string) (string, error) {
	fn, ok := reqFilterKeyFns[s]
	if ok {
		return fn(req)
	}

	return s, nil
}

func matchReqStringLiteral(req *http.Request, strLiteral filter.StringLiteral) (bool, error) {
	for key, values := range req.Header {
		for _, value := range values {
			if strings.Contains(
				strings.ToLower(fmt.Sprintf("%v: %v", key, value)),
				strings.ToLower(strLiteral.Value),
			) {
				return true, nil
			}
		}
	}

	for _, fn := range reqFilterKeyFns {
		value, err := fn(req)
		if err != nil {
			return false, err
		}

		if strings.Contains(strings.ToLower(value), strings.ToLower(strLiteral.Value)) {
			return true, nil
		}
	}

	return false, nil
}

func MatchRequestScope(req *http.Request, s *scope.Scope) (bool, error) {
	for _, rule := range s.Rules() {
		if rule.URL != nil && req.URL != nil {
			if matches := rule.URL.MatchString(req.URL.String()); matches {
				return true, nil
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
				return true, nil
			case rule.Header.Key == nil && rule.Header.Value != nil && valueMatches:
				return true, nil
			case rule.Header.Key != nil && rule.Header.Value != nil && keyMatches && valueMatches:
				return true, nil
			}
		}

		if rule.Body != nil {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return false, fmt.Errorf("failed to read request body: %w", err)
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			if matches := rule.Body.Match(body); matches {
				return true, nil
			}
		}
	}

	return false, nil
}

// MatchResponseFilter returns true if an HTTP response matches the response filter expression.
func MatchResponseFilter(res *http.Response, expr filter.Expression) (bool, error) {
	switch e := expr.(type) {
	case filter.PrefixExpression:
		return matchResPrefixExpr(res, e)
	case filter.InfixExpression:
		return matchResInfixExpr(res, e)
	case filter.StringLiteral:
		return matchResStringLiteral(res, e)
	default:
		return false, fmt.Errorf("expression type (%T) not supported", expr)
	}
}

func matchResPrefixExpr(res *http.Response, expr filter.PrefixExpression) (bool, error) {
	switch expr.Operator {
	case filter.TokOpNot:
		match, err := MatchResponseFilter(res, expr.Right)
		if err != nil {
			return false, err
		}

		return !match, nil
	default:
		return false, errors.New("operator is not supported")
	}
}

func matchResInfixExpr(res *http.Response, expr filter.InfixExpression) (bool, error) {
	switch expr.Operator {
	case filter.TokOpAnd:
		left, err := MatchResponseFilter(res, expr.Left)
		if err != nil {
			return false, err
		}

		right, err := MatchResponseFilter(res, expr.Right)
		if err != nil {
			return false, err
		}

		return left && right, nil
	case filter.TokOpOr:
		left, err := MatchResponseFilter(res, expr.Left)
		if err != nil {
			return false, err
		}

		right, err := MatchResponseFilter(res, expr.Right)
		if err != nil {
			return false, err
		}

		return left || right, nil
	}

	left, ok := expr.Left.(filter.StringLiteral)
	if !ok {
		return false, errors.New("left operand must be a string literal")
	}

	leftVal, err := getMappedStringLiteralFromRes(res, left.Value)
	if err != nil {
		return false, fmt.Errorf("failed to get string literal from response for left operand: %w", err)
	}

	if leftVal == "headers" {
		match, err := filter.MatchHTTPHeaders(expr.Operator, expr.Right, res.Header)
		if err != nil {
			return false, fmt.Errorf("failed to match request HTTP headers: %w", err)
		}

		return match, nil
	}

	if expr.Operator == filter.TokOpRe || expr.Operator == filter.TokOpNotRe {
		right, ok := expr.Right.(filter.RegexpLiteral)
		if !ok {
			return false, errors.New("right operand must be a regular expression")
		}

		switch expr.Operator {
		case filter.TokOpRe:
			return right.MatchString(leftVal), nil
		case filter.TokOpNotRe:
			return !right.MatchString(leftVal), nil
		}
	}

	right, ok := expr.Right.(filter.StringLiteral)
	if !ok {
		return false, errors.New("right operand must be a string literal")
	}

	rightVal, err := getMappedStringLiteralFromRes(res, right.Value)
	if err != nil {
		return false, fmt.Errorf("failed to get string literal from response for right operand: %w", err)
	}

	switch expr.Operator {
	case filter.TokOpEq:
		return leftVal == rightVal, nil
	case filter.TokOpNotEq:
		return leftVal != rightVal, nil
	case filter.TokOpGt:
		// TODO(?) attempt to parse as int.
		return leftVal > rightVal, nil
	case filter.TokOpLt:
		// TODO(?) attempt to parse as int.
		return leftVal < rightVal, nil
	case filter.TokOpGtEq:
		// TODO(?) attempt to parse as int.
		return leftVal >= rightVal, nil
	case filter.TokOpLtEq:
		// TODO(?) attempt to parse as int.
		return leftVal <= rightVal, nil
	default:
		return false, errors.New("unsupported operator")
	}
}

func getMappedStringLiteralFromRes(res *http.Response, s string) (string, error) {
	fn, ok := resFilterKeyFns[s]
	if ok {
		return fn(res)
	}

	return s, nil
}

func matchResStringLiteral(res *http.Response, strLiteral filter.StringLiteral) (bool, error) {
	for key, values := range res.Header {
		for _, value := range values {
			if strings.Contains(
				strings.ToLower(fmt.Sprintf("%v: %v", key, value)),
				strings.ToLower(strLiteral.Value),
			) {
				return true, nil
			}
		}
	}

	for _, fn := range resFilterKeyFns {
		value, err := fn(res)
		if err != nil {
			return false, err
		}

		if strings.Contains(strings.ToLower(value), strings.ToLower(strLiteral.Value)) {
			return true, nil
		}
	}

	return false, nil
}
