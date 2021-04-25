package search

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		expectedExpression Expression
		expectedError      error
	}{
		{
			name:               "empty query",
			input:              "",
			expectedExpression: nil,
			expectedError:      errors.New("search: unexpected EOF"),
		},
		{
			name:               "string literal expression",
			input:              "foobar",
			expectedExpression: &StringLiteral{Value: "foobar"},
			expectedError:      nil,
		},
		{
			name:  "boolean expression with equal operator",
			input: "foo = bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpEq,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with not equal operator",
			input: "foo != bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpNotEq,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with greater than operator",
			input: "foo > bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpGt,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with less than operator",
			input: "foo < bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpLt,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with greater than or equal operator",
			input: "foo >= bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpGtEq,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with less than or equal operator",
			input: "foo <= bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpLtEq,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with regular expression operator",
			input: "foo =~ bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpRe,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with not regular expression operator",
			input: "foo !~ bar",
			expectedExpression: &InfixExpression{
				Operator: TokOpNotRe,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with AND, OR and NOT operators",
			input: "foo AND bar OR NOT baz",
			expectedExpression: &InfixExpression{
				Operator: TokOpAnd,
				Left:     &StringLiteral{Value: "foo"},
				Right: &InfixExpression{
					Operator: TokOpOr,
					Left:     &StringLiteral{Value: "bar"},
					Right: &PrefixExpression{
						Operator: TokOpNot,
						Right:    &StringLiteral{Value: "baz"},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:  "boolean expression with nested group",
			input: "(foo AND bar) OR NOT baz",
			expectedExpression: &InfixExpression{
				Operator: TokOpOr,
				Left: &InfixExpression{
					Operator: TokOpAnd,
					Left:     &StringLiteral{Value: "foo"},
					Right:    &StringLiteral{Value: "bar"},
				},
				Right: &PrefixExpression{
					Operator: TokOpNot,
					Right:    &StringLiteral{Value: "baz"},
				},
			},
			expectedError: nil,
		},
		{
			name:  "implicit boolean expression with string literal operands",
			input: "foo bar baz",
			expectedExpression: &InfixExpression{
				Operator: TokOpAnd,
				Left: &InfixExpression{
					Operator: TokOpAnd,
					Left:     &StringLiteral{Value: "foo"},
					Right:    &StringLiteral{Value: "bar"},
				},
				Right: &StringLiteral{Value: "baz"},
			},
			expectedError: nil,
		},
		{
			name:  "implicit boolean expression nested in group",
			input: "(foo bar)",
			expectedExpression: &InfixExpression{
				Operator: TokOpAnd,
				Left:     &StringLiteral{Value: "foo"},
				Right:    &StringLiteral{Value: "bar"},
			},
			expectedError: nil,
		},
		{
			name:  "implicit and explicit boolean expression with string literal operands",
			input: "foo bar OR baz yolo",
			expectedExpression: &InfixExpression{
				Operator: TokOpAnd,
				Left: &InfixExpression{
					Operator: TokOpAnd,
					Left:     &StringLiteral{Value: "foo"},
					Right: &InfixExpression{
						Operator: TokOpOr,
						Left:     &StringLiteral{Value: "bar"},
						Right:    &StringLiteral{Value: "baz"},
					},
				},
				Right: &StringLiteral{Value: "yolo"},
			},
			expectedError: nil,
		},
		{
			name:  "implicit boolean expression with comparison operands",
			input: "foo=bar baz=~yolo",
			expectedExpression: &InfixExpression{
				Operator: TokOpAnd,
				Left: &InfixExpression{
					Operator: TokOpEq,
					Left:     &StringLiteral{Value: "foo"},
					Right:    &StringLiteral{Value: "bar"},
				},
				Right: &InfixExpression{
					Operator: TokOpRe,
					Left:     &StringLiteral{Value: "baz"},
					Right:    &StringLiteral{Value: "yolo"},
				},
			},
			expectedError: nil,
		},
		{
			name:  "eq operator takes precedence over boolean ops",
			input: "foo=bar OR baz=yolo",
			expectedExpression: &InfixExpression{
				Operator: TokOpOr,
				Left: &InfixExpression{
					Operator: TokOpEq,
					Left:     &StringLiteral{Value: "foo"},
					Right:    &StringLiteral{Value: "bar"},
				},
				Right: &InfixExpression{
					Operator: TokOpEq,
					Left:     &StringLiteral{Value: "baz"},
					Right:    &StringLiteral{Value: "yolo"},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseQuery(tt.input)
			assertError(t, tt.expectedError, err)
			if !reflect.DeepEqual(tt.expectedExpression, got) {
				t.Errorf("expected: %v, got: %v", tt.expectedExpression, got)
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
