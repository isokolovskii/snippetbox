package main

import (
	"net/http"
	"testing"

	"snippetbox.isokol.dev/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, healthCheckRoute)
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, ok)
}
