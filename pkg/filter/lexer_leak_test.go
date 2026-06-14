package filter

import (
	"runtime"
	"testing"
	"time"
)

// TestParseQueryNoGoroutineLeak is a regression test for a goroutine leak in
// the lexer. NewLexer spawns a goroutine that emits tokens on an unbuffered
// channel. When ParseQuery returns early on a parse error, it used to stop
// reading tokens, leaving that goroutine blocked forever on a channel send.
// Because ParseQuery is reachable from the (unauthenticated) GraphQL API via
// search expressions, repeated malformed queries leaked one goroutine each,
// enabling a memory/goroutine-exhaustion denial of service.
func TestParseQueryNoGoroutineLeak(t *testing.T) {
	// Inputs that make ParseQuery return an error before the lexer reaches EOF.
	malformed := []string{
		"= foo",
		"!= bar",
		"(foo",
		"=~ baz",
		"NOT",
	}

	runtime.GC()
	time.Sleep(50 * time.Millisecond)
	before := runtime.NumGoroutine()

	const iterations = 500
	for i := 0; i < iterations; i++ {
		input := malformed[i%len(malformed)]
		if _, err := ParseQuery(input); err == nil {
			t.Fatalf("expected parse error for input %q", input)
		}
	}

	// Give any abandoned goroutines a chance to settle.
	deadline := time.Now().Add(2 * time.Second)
	var after int
	for time.Now().Before(deadline) {
		runtime.GC()
		after = runtime.NumGoroutine()
		if after-before <= 5 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if leaked := after - before; leaked > 5 {
		t.Fatalf("goroutine leak: %d goroutines still running after %d malformed queries (before=%d, after=%d)",
			leaked, iterations, before, after)
	}
}
