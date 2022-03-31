package search

import (
	"errors"
	"fmt"
	"net/http"
)

func MatchHTTPHeaders(op TokenType, expr Expression, headers http.Header) (bool, error) {
	if headers == nil {
		return false, nil
	}

	switch op {
	case TokOpEq:
		strLiteral, ok := expr.(StringLiteral)
		if !ok {
			return false, errors.New("search: expression must be a string literal")
		}

		// Return `true` if at least one header (<key>: <value>) is equal to the string literal.
		for key, values := range headers {
			for _, value := range values {
				if strLiteral.Value == fmt.Sprintf("%v: %v", key, value) {
					return true, nil
				}
			}
		}

		return false, nil
	case TokOpNotEq:
		strLiteral, ok := expr.(StringLiteral)
		if !ok {
			return false, errors.New("search: expression must be a string literal")
		}

		// Return `true` if none of the headers (<key>: <value>) are equal to the string literal.
		for key, values := range headers {
			for _, value := range values {
				if strLiteral.Value == fmt.Sprintf("%v: %v", key, value) {
					return false, nil
				}
			}
		}

		return true, nil
	case TokOpRe:
		re, ok := expr.(RegexpLiteral)
		if !ok {
			return false, errors.New("search: expression must be a regular expression")
		}

		// Return `true` if at least one header (<key>: <value>) matches the regular expression.
		for key, values := range headers {
			for _, value := range values {
				if re.MatchString(fmt.Sprintf("%v: %v", key, value)) {
					return true, nil
				}
			}
		}

		return false, nil
	case TokOpNotRe:
		re, ok := expr.(RegexpLiteral)
		if !ok {
			return false, errors.New("search: expression must be a regular expression")
		}

		// Return `true` if none of the headers (<key>: <value>) match the regular expression.
		for key, values := range headers {
			for _, value := range values {
				if re.MatchString(fmt.Sprintf("%v: %v", key, value)) {
					return false, nil
				}
			}
		}

		return true, nil
	default:
		return false, fmt.Errorf("search: unsupported operator %q", op.String())
	}
}
