package sqlite

import (
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/dstotijn/hetty/pkg/search"
)

func TestParseSearchExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		searchExpr      search.Expression
		expectedSqlizer sq.Sqlizer
		expectedError   error
	}{
		{
			name: "req.body = bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpEq,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.Eq{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body != bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpNotEq,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.NotEq{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body > bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpGt,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.Gt{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body < bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpLt,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.Lt{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body >= bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpGtEq,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.GtOrEq{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body <= bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpLtEq,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.LtOrEq{"req.body": "bar"},
			expectedError:   nil,
		},
		{
			name: "req.body =~ bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpRe,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.Expr("regexp(?, req.body)", "bar"),
			expectedError:   nil,
		},
		{
			name: "req.body !~ bar",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpNotRe,
				Left:     &search.StringLiteral{Value: "req.body"},
				Right:    &search.StringLiteral{Value: "bar"},
			},
			expectedSqlizer: sq.Expr("NOT regexp(?, req.body)", "bar"),
			expectedError:   nil,
		},
		{
			name: "req.body = bar AND res.body = yolo",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpAnd,
				Left: &search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     &search.StringLiteral{Value: "req.body"},
					Right:    &search.StringLiteral{Value: "bar"},
				},
				Right: &search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     &search.StringLiteral{Value: "res.body"},
					Right:    &search.StringLiteral{Value: "yolo"},
				},
			},
			expectedSqlizer: sq.And{
				sq.Eq{"req.body": "bar"},
				sq.Eq{"res.body": "yolo"},
			},
			expectedError: nil,
		},
		{
			name: "req.body = bar AND res.body = yolo AND req.method = POST",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpAnd,
				Left: &search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     &search.StringLiteral{Value: "req.body"},
					Right:    &search.StringLiteral{Value: "bar"},
				},
				Right: &search.InfixExpression{
					Operator: search.TokOpAnd,
					Left: &search.InfixExpression{
						Operator: search.TokOpEq,
						Left:     &search.StringLiteral{Value: "res.body"},
						Right:    &search.StringLiteral{Value: "yolo"},
					},
					Right: &search.InfixExpression{
						Operator: search.TokOpEq,
						Left:     &search.StringLiteral{Value: "req.method"},
						Right:    &search.StringLiteral{Value: "POST"},
					},
				},
			},
			expectedSqlizer: sq.And{
				sq.Eq{"req.body": "bar"},
				sq.And{
					sq.Eq{"res.body": "yolo"},
					sq.Eq{"req.method": "POST"},
				},
			},
			expectedError: nil,
		},
		{
			name: "req.body = bar OR res.body = yolo",
			searchExpr: &search.InfixExpression{
				Operator: search.TokOpOr,
				Left: &search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     &search.StringLiteral{Value: "req.body"},
					Right:    &search.StringLiteral{Value: "bar"},
				},
				Right: &search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     &search.StringLiteral{Value: "res.body"},
					Right:    &search.StringLiteral{Value: "yolo"},
				},
			},
			expectedSqlizer: sq.Or{
				sq.Eq{"req.body": "bar"},
				sq.Eq{"res.body": "yolo"},
			},
			expectedError: nil,
		},
		{
			name: "foo",
			searchExpr: &search.StringLiteral{
				Value: "foo",
			},
			expectedSqlizer: sq.Or{
				sq.Like{"req.body": "%foo%"},
				sq.Like{"req.id": "%foo%"},
				sq.Like{"req.method": "%foo%"},
				sq.Like{"req.proto": "%foo%"},
				sq.Like{"req.timestamp": "%foo%"},
				sq.Like{"req.url": "%foo%"},
				sq.Like{"res.body": "%foo%"},
				sq.Like{"res.id": "%foo%"},
				sq.Like{"res.proto": "%foo%"},
				sq.Like{"res.status_code": "%foo%"},
				sq.Like{"res.status_reason": "%foo%"},
				sq.Like{"res.timestamp": "%foo%"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseSearchExpr(tt.searchExpr)
			assertError(t, tt.expectedError, err)
			if !reflect.DeepEqual(tt.expectedSqlizer, got) {
				t.Errorf("expected: %#v, got: %#v", tt.expectedSqlizer, got)
			}
		})
	}
}

func assertError(t *testing.T, exp, got error) {
	t.Helper()

	switch {
	case exp == nil && got != nil:
		t.Fatalf("expected: nil, got: %v", got)
	case exp != nil && got == nil:
		t.Fatalf("expected: %v, got: nil", exp.Error())
	case exp != nil && got != nil && exp.Error() != got.Error():
		t.Fatalf("expected: %v, got: %v", exp.Error(), got.Error())
	}
}
