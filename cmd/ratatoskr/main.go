// Command ratatoskr is the Ratatoskr web security proxy.
//
// It currently only prints its version; wiring lands in later tasks
// (see docs/08-roadmap.md).
package main

import (
	"fmt"
	"os"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if _, err := fmt.Fprintf(os.Stdout, "ratatoskr %s\n", version); err != nil {
		os.Exit(1)
	}
}
