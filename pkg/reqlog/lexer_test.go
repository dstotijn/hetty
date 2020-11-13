package reqlog

import (
	"strings"
	"testing"

	"github.com/db47h/lex"
)

func TestLex(t *testing.T) {
	lexTests := []struct {
		name     string
		input    string
		expected []lexItem
	}{
		{
			name:  "empty query",
			input: "",
			expected: []lexItem{
				{tokEOF, ""},
			},
		},
		{
			name:  "single unquoted value",
			input: "foobar",
			expected: []lexItem{
				{tokString, "foobar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "single unquoted value with non letter",
			input: "foob*",
			expected: []lexItem{
				{tokString, "foob*"},
				{tokEOF, ""},
			},
		},
		{
			name:  "multiple unquoted values",
			input: "foo bar",
			expected: []lexItem{
				{tokString, "foo"},
				{tokString, "bar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "quoted value",
			input: `"foo bar"`,
			expected: []lexItem{
				{tokString, "foo bar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with negation operator",
			input: "NOT foobar",
			expected: []lexItem{
				{tokOpNot, ""},
				{tokString, "foobar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with and operator",
			input: "foo AND bar",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpAnd, ""},
				{tokString, "bar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with or operator",
			input: "foo OR bar",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpOr, ""},
				{tokString, "bar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with equals operator",
			input: "foo = bar",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpEq, ""},
				{tokString, "bar"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with greater than operator",
			input: "foo > 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpGt, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with greater than or equal operator",
			input: "foo >= 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpGteq, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with less than operator",
			input: "foo < 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpLt, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with less than or equal operator",
			input: "foo <= 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpLteq, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with regular expression operator",
			input: "foo =~ 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpRe, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with not regular expression operator",
			input: "foo !~ 42",
			expected: []lexItem{
				{tokString, "foo"},
				{tokOpNre, ""},
				{tokString, "42"},
				{tokEOF, ""},
			},
		},
		{
			name:  "comparison with parentheses",
			input: "(foo OR bar) AND baz",
			expected: []lexItem{
				{tokParenOpen, ""},
				{tokString, "foo"},
				{tokOpOr, ""},
				{tokString, "bar"},
				{tokParenClose, ""},
				{tokOpAnd, ""},
				{tokString, "baz"},
				{tokEOF, ""},
			},
		},
	}

	for _, tt := range lexTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			file := lex.NewFile(tt.name, strings.NewReader(tt.input))
			l := lex.NewLexer(file, lexQuery)

			for i, exp := range tt.expected {
				token, _, value := l.Lex()
				if err, isErr := value.(error); isErr {
					t.Fatalf("unexpected error: %v", err)
				}
				valueStr, _ := value.(string)
				got := lexItem{
					token: token,
					value: valueStr,
				}
				if got != exp {
					t.Errorf("%v: got: %+v, expected: %+v", i, got, exp)
				}
			}
		})
	}
}
