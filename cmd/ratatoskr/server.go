package main

import (
	"io"
	"log/slog"
	"net/http"
)

// placeholderPage is served at "/" until the Svelte SPA lands (T-003).
const placeholderPage = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Ratatoskr</title>
  </head>
  <body>
    <h1>Ratatoskr</h1>
    <p>Placeholder page. The web UI is not built yet.</p>
  </body>
</html>
`

// newHandler returns the UI/API handler. Routes are added as the API lands; for
// now only the placeholder page exists.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "text/html; charset=utf-8")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Cache-Control", "no-store")

		if _, err := io.WriteString(w, placeholderPage); err != nil {
			slog.Debug("write placeholder page", "error", err)
		}
	})

	return mux
}
