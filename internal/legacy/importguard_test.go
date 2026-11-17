package legacy

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestNoLegacyImports fails when a package outside internal/legacy imports the
// reference tree. docs/03-architecture.md §3 keeps these packages as reading
// material only: importing one would put dead Hetty code back on a live path.
func TestNoLegacyImports(t *testing.T) {
	t.Parallel()

	pkgs := listModulePackages(t)

	modulePath := ""
	for _, p := range pkgs {
		if p.Module != nil && p.Module.Main {
			modulePath = p.Module.Path
			break
		}
	}
	if modulePath == "" {
		t.Fatal("go list reported no package of the main module")
	}

	legacyPath := modulePath + "/internal/legacy"
	isLegacy := func(path string) bool {
		return path == legacyPath || strings.HasPrefix(path, legacyPath+"/")
	}

	var offenders []string
	checked := 0

	for _, p := range pkgs {
		if p.Module == nil || !p.Module.Main || isLegacy(p.ImportPath) {
			// Stdlib and third-party packages cannot reach our internal/ tree.
			continue
		}

		// `go list -test` also reports synthetic variants such as
		// "example.com/foo [example.com/foo.test]" and "example.com/foo.test".
		// Only the real package path matters; its test imports are listed in
		// Imports as well.
		importPath, _, _ := strings.Cut(p.ImportPath, " [")
		if isLegacy(importPath) || strings.HasSuffix(importPath, ".test") {
			continue
		}

		checked++

		for _, imp := range p.Imports {
			if isLegacy(imp) {
				offenders = append(offenders, fmt.Sprintf("%s imports %s", importPath, imp))
			}
		}
	}

	if checked == 0 {
		t.Fatal("go list produced no package to inspect, the guard checked nothing")
	}

	slices.Sort(offenders)

	if offenders = slices.Compact(offenders); len(offenders) > 0 {
		t.Fatalf("internal/legacy is reference-only (see internal/legacy/README.md), but:\n%s",
			strings.Join(offenders, "\n"))
	}
}

// goPackage is the subset of `go list -json` output the guard needs.
type goPackage struct {
	ImportPath string
	Imports    []string
	Module     *struct {
		Path string
		Main bool
	}
}

// listModulePackages returns every package in the build graph of ./... (test
// dependencies included), with the go command rooted at the module directory.
func listModulePackages(t *testing.T) []goPackage {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "go", "list", "-deps", "-json", "-test", "./...")
	cmd.Dir = moduleRoot(t)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("go list: stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("go list: start: %v", err)
	}

	var pkgs []goPackage

	dec := jsontext.NewDecoder(stdout)
	for {
		var p goPackage

		err := json.UnmarshalDecode(dec, &p)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()

			t.Fatalf("go list: decode output: %v", err)
		}

		pkgs = append(pkgs, p)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatalf("go list: %v: %s", err, stderr.String())
	}

	return pkgs
}

// moduleRoot walks up from the test's working directory to the module root.
func moduleRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("no go.mod above %s", dir)
		}

		dir = parent
	}
}
