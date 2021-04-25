package sqlite

import (
	"errors"
	"fmt"
	"sort"

	sq "github.com/Masterminds/squirrel"

	"github.com/dstotijn/hetty/pkg/search"
)

var stringLiteralMap = map[string]string{
	// http_requests
	"req.id":        "req.id",
	"req.proto":     "req.proto",
	"req.url":       "req.url",
	"req.method":    "req.method",
	"req.body":      "req.body",
	"req.timestamp": "req.timestamp",
	// http_responses
	"res.id":           "res.id",
	"res.proto":        "res.proto",
	"res.statusCode":   "res.status_code",
	"res.statusReason": "res.status_reason",
	"res.body":         "res.body",
	"res.timestamp":    "res.timestamp",
	// TODO: http_headers
}

func parseSearchExpr(expr search.Expression) (sq.Sqlizer, error) {
	switch e := expr.(type) {
	case *search.PrefixExpression:
		return parsePrefixExpr(e)
	case *search.InfixExpression:
		return parseInfixExpr(e)
	case *search.StringLiteral:
		return parseStringLiteral(e)
	default:
		return nil, fmt.Errorf("expression type (%v) not supported", expr)
	}
}

func parsePrefixExpr(expr *search.PrefixExpression) (sq.Sqlizer, error) {
	switch expr.Operator {
	case search.TokOpNot:
		// TODO: Find a way to prefix an `sq.Sqlizer` with "NOT".
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("operator is not supported")
	}
}

func parseInfixExpr(expr *search.InfixExpression) (sq.Sqlizer, error) {
	switch expr.Operator {
	case search.TokOpAnd:
		left, err := parseSearchExpr(expr.Left)
		if err != nil {
			return nil, err
		}

		right, err := parseSearchExpr(expr.Right)
		if err != nil {
			return nil, err
		}

		return sq.And{left, right}, nil
	case search.TokOpOr:
		left, err := parseSearchExpr(expr.Left)
		if err != nil {
			return nil, err
		}

		right, err := parseSearchExpr(expr.Right)
		if err != nil {
			return nil, err
		}

		return sq.Or{left, right}, nil
	}

	left, ok := expr.Left.(*search.StringLiteral)
	if !ok {
		return nil, errors.New("left operand must be a string literal")
	}

	right, ok := expr.Right.(*search.StringLiteral)
	if !ok {
		return nil, errors.New("right operand must be a string literal")
	}

	mappedLeft, ok := stringLiteralMap[left.Value]
	if !ok {
		return nil, fmt.Errorf("invalid string literal: %v", left)
	}

	switch expr.Operator {
	case search.TokOpEq:
		return sq.Eq{mappedLeft: right.Value}, nil
	case search.TokOpNotEq:
		return sq.NotEq{mappedLeft: right.Value}, nil
	case search.TokOpGt:
		return sq.Gt{mappedLeft: right.Value}, nil
	case search.TokOpLt:
		return sq.Lt{mappedLeft: right.Value}, nil
	case search.TokOpGtEq:
		return sq.GtOrEq{mappedLeft: right.Value}, nil
	case search.TokOpLtEq:
		return sq.LtOrEq{mappedLeft: right.Value}, nil
	case search.TokOpRe:
		return sq.Expr(fmt.Sprintf("regexp(?, %v)", mappedLeft), right.Value), nil
	case search.TokOpNotRe:
		return sq.Expr(fmt.Sprintf("NOT regexp(?, %v)", mappedLeft), right.Value), nil
	default:
		return nil, errors.New("unsupported operator")
	}
}

func parseStringLiteral(strLiteral *search.StringLiteral) (sq.Sqlizer, error) {
	// Sorting is not necessary, but makes it easier to do assertions in tests.
	sortedKeys := make([]string, 0, len(stringLiteralMap))

	for _, v := range stringLiteralMap {
		sortedKeys = append(sortedKeys, v)
	}

	sort.Strings(sortedKeys)

	or := make(sq.Or, len(stringLiteralMap))
	for i, value := range sortedKeys {
		or[i] = sq.Like{value: "%" + strLiteral.Value + "%"}
	}

	return or, nil
}
