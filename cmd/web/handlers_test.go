package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.isokol.dev/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	app := application{
		logger: slog.New(slog.DiscardHandler),
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL+healthCheckRoute, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), ok)
}
