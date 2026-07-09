package filter

import (
	"errors"
	"strings"
	"testing"
)

// TestParseQueryDepthLimit is a regression test for an unrecoverable stack
// overflow in the recursive-descent parser. Deeply nested input (many "(" or
// repeated "NOT") used to recurse until the goroutine stack was exhausted,
// which is a fatal runtime error in Go that recover() cannot catch and that
// therefore crashes the whole process. ParseQuery is reachable from the
// unauthenticated GraphQL API, making this a remote denial of service.
func TestParseQueryDepthLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"deeply nested parentheses", strings.Repeat("(", 1_000_000)},
		{"repeated NOT prefix", strings.Repeat("NOT ", 1_000_000)},
		{"deeply nested groups", strings.Repeat("(", 100_000) + "foo" + strings.Repeat(")", 100_000)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Must return an error instead of crashing the process.
			_, err := ParseQuery(tt.input)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !errors.Is(err, ErrMaxDepthExceeded) {
				t.Fatalf("expected ErrMaxDepthExceeded, got: %v", err)
			}
		})
	}
}

// TestParseQueryReasonableNestingStillWorks ensures the depth guard doesn't
// reject legitimate, modestly nested expressions.
func TestParseQueryReasonableNestingStillWorks(t *testing.T) {
	t.Parallel()

	input := strings.Repeat("(", 10) + "foo" + strings.Repeat(")", 10)
	if _, err := ParseQuery(input); err != nil {
		t.Fatalf("legitimate nested expression should parse, got error: %v", err)
	}
}
