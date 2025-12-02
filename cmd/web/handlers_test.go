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

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid ID",
			urlPath:    snippetViewRoute + "/1",
			wantStatus: http.StatusOK,
			wantBody:   "Mock snippet",
		},
		{
			name:       "Non-existent ID",
			urlPath:    snippetViewRoute + "/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Negative ID",
			urlPath:    snippetViewRoute + "/-1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Decimal ID",
			urlPath:    snippetViewRoute + "/1.23",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "String ID",
			urlPath:    snippetViewRoute + "/foo",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Empty ID",
			urlPath:    snippetViewRoute + "/",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, tt.urlPath)
			assert.Equal(t, res.status, tt.wantStatus)
			assert.StringContains(t, res.body, tt.wantBody)
		})
	}
}
