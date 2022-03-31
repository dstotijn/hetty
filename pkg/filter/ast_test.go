package filter_test

import (
	"regexp"
	"testing"

	"github.com/dstotijn/hetty/pkg/filter"
)

func TestExpressionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expression filter.Expression
		expected   string
	}{
		{
			name:       "string literal expression",
			expression: filter.StringLiteral{Value: "foobar"},
			expected:   `"foobar"`,
		},
		{
			name: "boolean expression with equal operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpEq,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" = "bar")`,
		},
		{
			name: "boolean expression with not equal operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpNotEq,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" != "bar")`,
		},
		{
			name: "boolean expression with greater than operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpGt,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" > "bar")`,
		},
		{
			name: "boolean expression with less than operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpLt,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" < "bar")`,
		},
		{
			name: "boolean expression with greater than or equal operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpGtEq,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" >= "bar")`,
		},
		{
			name: "boolean expression with less than or equal operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpLtEq,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" <= "bar")`,
		},
		{
			name: "boolean expression with regular expression operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpRe,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.RegexpLiteral{regexp.MustCompile("bar")},
			},
			expected: `("foo" =~ "bar")`,
		},
		{
			name: "boolean expression with not regular expression operator",
			expression: filter.InfixExpression{
				Operator: filter.TokOpNotRe,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.RegexpLiteral{regexp.MustCompile("bar")},
			},
			expected: `("foo" !~ "bar")`,
		},
		{
			name: "boolean expression with AND, OR and NOT operators",
			expression: filter.InfixExpression{
				Operator: filter.TokOpAnd,
				Left:     filter.StringLiteral{Value: "foo"},
				Right: filter.InfixExpression{
					Operator: filter.TokOpOr,
					Left:     filter.StringLiteral{Value: "bar"},
					Right: filter.PrefixExpression{
						Operator: filter.TokOpNot,
						Right:    filter.StringLiteral{Value: "baz"},
					},
				},
			},
			expected: `("foo" AND ("bar" OR (NOT "baz")))`,
		},
		{
			name: "boolean expression with nested group",
			expression: filter.InfixExpression{
				Operator: filter.TokOpOr,
				Left: filter.InfixExpression{
					Operator: filter.TokOpAnd,
					Left:     filter.StringLiteral{Value: "foo"},
					Right:    filter.StringLiteral{Value: "bar"},
				},
				Right: filter.PrefixExpression{
					Operator: filter.TokOpNot,
					Right:    filter.StringLiteral{Value: "baz"},
				},
			},
			expected: `(("foo" AND "bar") OR (NOT "baz"))`,
		},
		{
			name: "implicit boolean expression with string literal operands",
			expression: filter.InfixExpression{
				Operator: filter.TokOpAnd,
				Left: filter.InfixExpression{
					Operator: filter.TokOpAnd,
					Left:     filter.StringLiteral{Value: "foo"},
					Right:    filter.StringLiteral{Value: "bar"},
				},
				Right: filter.StringLiteral{Value: "baz"},
			},
			expected: `(("foo" AND "bar") AND "baz")`,
		},
		{
			name: "implicit boolean expression nested in group",
			expression: filter.InfixExpression{
				Operator: filter.TokOpAnd,
				Left:     filter.StringLiteral{Value: "foo"},
				Right:    filter.StringLiteral{Value: "bar"},
			},
			expected: `("foo" AND "bar")`,
		},
		{
			name: "implicit and explicit boolean expression with string literal operands",
			expression: filter.InfixExpression{
				Operator: filter.TokOpAnd,
				Left: filter.InfixExpression{
					Operator: filter.TokOpAnd,
					Left:     filter.StringLiteral{Value: "foo"},
					Right: filter.InfixExpression{
						Operator: filter.TokOpOr,
						Left:     filter.StringLiteral{Value: "bar"},
						Right:    filter.StringLiteral{Value: "baz"},
					},
				},
				Right: filter.StringLiteral{Value: "yolo"},
			},
			expected: `(("foo" AND ("bar" OR "baz")) AND "yolo")`,
		},
		{
			name: "implicit boolean expression with comparison operands",
			expression: filter.InfixExpression{
				Operator: filter.TokOpAnd,
				Left: filter.InfixExpression{
					Operator: filter.TokOpEq,
					Left:     filter.StringLiteral{Value: "foo"},
					Right:    filter.StringLiteral{Value: "bar"},
				},
				Right: filter.InfixExpression{
					Operator: filter.TokOpRe,
					Left:     filter.StringLiteral{Value: "baz"},
					Right:    filter.RegexpLiteral{regexp.MustCompile("yolo")},
				},
			},
			expected: `(("foo" = "bar") AND ("baz" =~ "yolo"))`,
		},
		{
			name: "eq operator takes precedence over boolean ops",
			expression: filter.InfixExpression{
				Operator: filter.TokOpOr,
				Left: filter.InfixExpression{
					Operator: filter.TokOpEq,
					Left:     filter.StringLiteral{Value: "foo"},
					Right:    filter.StringLiteral{Value: "bar"},
				},
				Right: filter.InfixExpression{
					Operator: filter.TokOpEq,
					Left:     filter.StringLiteral{Value: "baz"},
					Right:    filter.StringLiteral{Value: "yolo"},
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
