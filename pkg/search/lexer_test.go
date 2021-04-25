package search

import "testing"

func TestNextToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "unquoted string",
			input: "foo bar",
			expected: []Token{
				{TokString, "foo"},
				{TokString, "bar"},
				{TokEOF, ""},
			},
		},
		{
			name:  "quoted string",
			input: `"foo bar" "baz"`,
			expected: []Token{
				{TokString, "foo bar"},
				{TokString, "baz"},
				{TokEOF, ""},
			},
		},
		{
			name:  "boolean operator token types",
			input: "NOT AND OR",
			expected: []Token{
				{TokOpNot, "NOT"},
				{TokOpAnd, "AND"},
				{TokOpOr, "OR"},
				{TokEOF, ""},
			},
		},
		{
			name:  "comparison operator token types",
			input: `= != < > <= >= =~ !~`,
			expected: []Token{
				{TokOpEq, "="},
				{TokOpNotEq, "!="},
				{TokOpLt, "<"},
				{TokOpGt, ">"},
				{TokOpLtEq, "<="},
				{TokOpGtEq, ">="},
				{TokOpRe, "=~"},
				{TokOpNotRe, "!~"},
				{TokEOF, ""},
			},
		},
		{
			name:  "with parentheses",
			input: "(foo AND bar) OR baz",
			expected: []Token{
				{TokParenOpen, "("},
				{TokString, "foo"},
				{TokOpAnd, "AND"},
				{TokString, "bar"},
				{TokParenClose, ")"},
				{TokOpOr, "OR"},
				{TokString, "baz"},
				{TokEOF, ""},
			},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := NewLexer(tt.input)

			for _, exp := range tt.expected {
				got := l.Next()
				if got.Type != exp.Type {
					t.Errorf("invalid type (idx: %v, expected: %v, got: %v)",
						i, exp.Type, got.Type)
				}
				if got.Literal != exp.Literal {
					t.Errorf("invalid literal (idx: %v, expected: %v, got: %v)",
						i, exp.Literal, got.Literal)
				}
			}
		})
	}
}
