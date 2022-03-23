package search_test

import (
	"regexp"
	"testing"

	"github.com/dstotijn/hetty/pkg/search"
)

func TestExpressionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expression search.Expression
		expected   string
	}{
		{
			name:       "string literal expression",
			expression: search.StringLiteral{Value: "foobar"},
			expected:   `"foobar"`,
		},
		{
			name: "boolean expression with equal operator",
			expression: search.InfixExpression{
				Operator: search.TokOpEq,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" = "bar")`,
		},
		{
			name: "boolean expression with not equal operator",
			expression: search.InfixExpression{
				Operator: search.TokOpNotEq,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" != "bar")`,
		},
		{
			name: "boolean expression with greater than operator",
			expression: search.InfixExpression{
				Operator: search.TokOpGt,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" > "bar")`,
		},
		{
			name: "boolean expression with less than operator",
			expression: search.InfixExpression{
				Operator: search.TokOpLt,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" < "bar")`,
		},
		{
			name: "boolean expression with greater than or equal operator",
			expression: search.InfixExpression{
				Operator: search.TokOpGtEq,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" >= "bar")`,
		},
		{
			name: "boolean expression with less than or equal operator",
			expression: search.InfixExpression{
				Operator: search.TokOpLtEq,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" <= "bar")`,
		},
		{
			name: "boolean expression with regular expression operator",
			expression: search.InfixExpression{
				Operator: search.TokOpRe,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.RegexpLiteral{regexp.MustCompile("bar")},
			},
			expected: `("foo" =~ "bar")`,
		},
		{
			name: "boolean expression with not regular expression operator",
			expression: search.InfixExpression{
				Operator: search.TokOpNotRe,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.RegexpLiteral{regexp.MustCompile("bar")},
			},
			expected: `("foo" !~ "bar")`,
		},
		{
			name: "boolean expression with AND, OR and NOT operators",
			expression: search.InfixExpression{
				Operator: search.TokOpAnd,
				Left:     search.StringLiteral{Value: "foo"},
				Right: search.InfixExpression{
					Operator: search.TokOpOr,
					Left:     search.StringLiteral{Value: "bar"},
					Right: search.PrefixExpression{
						Operator: search.TokOpNot,
						Right:    search.StringLiteral{Value: "baz"},
					},
				},
			},
			expected: `("foo" AND ("bar" OR (NOT "baz")))`,
		},
		{
			name: "boolean expression with nested group",
			expression: search.InfixExpression{
				Operator: search.TokOpOr,
				Left: search.InfixExpression{
					Operator: search.TokOpAnd,
					Left:     search.StringLiteral{Value: "foo"},
					Right:    search.StringLiteral{Value: "bar"},
				},
				Right: search.PrefixExpression{
					Operator: search.TokOpNot,
					Right:    search.StringLiteral{Value: "baz"},
				},
			},
			expected: `(("foo" AND "bar") OR (NOT "baz"))`,
		},
		{
			name: "implicit boolean expression with string literal operands",
			expression: search.InfixExpression{
				Operator: search.TokOpAnd,
				Left: search.InfixExpression{
					Operator: search.TokOpAnd,
					Left:     search.StringLiteral{Value: "foo"},
					Right:    search.StringLiteral{Value: "bar"},
				},
				Right: search.StringLiteral{Value: "baz"},
			},
			expected: `(("foo" AND "bar") AND "baz")`,
		},
		{
			name: "implicit boolean expression nested in group",
			expression: search.InfixExpression{
				Operator: search.TokOpAnd,
				Left:     search.StringLiteral{Value: "foo"},
				Right:    search.StringLiteral{Value: "bar"},
			},
			expected: `("foo" AND "bar")`,
		},
		{
			name: "implicit and explicit boolean expression with string literal operands",
			expression: search.InfixExpression{
				Operator: search.TokOpAnd,
				Left: search.InfixExpression{
					Operator: search.TokOpAnd,
					Left:     search.StringLiteral{Value: "foo"},
					Right: search.InfixExpression{
						Operator: search.TokOpOr,
						Left:     search.StringLiteral{Value: "bar"},
						Right:    search.StringLiteral{Value: "baz"},
					},
				},
				Right: search.StringLiteral{Value: "yolo"},
			},
			expected: `(("foo" AND ("bar" OR "baz")) AND "yolo")`,
		},
		{
			name: "implicit boolean expression with comparison operands",
			expression: search.InfixExpression{
				Operator: search.TokOpAnd,
				Left: search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     search.StringLiteral{Value: "foo"},
					Right:    search.StringLiteral{Value: "bar"},
				},
				Right: search.InfixExpression{
					Operator: search.TokOpRe,
					Left:     search.StringLiteral{Value: "baz"},
					Right:    search.RegexpLiteral{regexp.MustCompile("yolo")},
				},
			},
			expected: `(("foo" = "bar") AND ("baz" =~ "yolo"))`,
		},
		{
			name: "eq operator takes precedence over boolean ops",
			expression: search.InfixExpression{
				Operator: search.TokOpOr,
				Left: search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     search.StringLiteral{Value: "foo"},
					Right:    search.StringLiteral{Value: "bar"},
				},
				Right: search.InfixExpression{
					Operator: search.TokOpEq,
					Left:     search.StringLiteral{Value: "baz"},
					Right:    search.StringLiteral{Value: "yolo"},
				},
			},
			expected: `(("foo" = "bar") OR ("baz" = "yolo"))`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.expression.String()
			if tt.expected != got {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}
