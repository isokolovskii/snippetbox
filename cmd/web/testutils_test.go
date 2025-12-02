package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"snippetbox.isokol.dev/internal/repositories"
	"snippetbox.isokol.dev/internal/repositories/mocks"
)

type (
	// Test http server.
	testServer struct {
		*httptest.Server
	}
	// Response of test http server.
	testResponse struct {
		// Response status code.
		status int
		// Response headers.
		headers http.Header
		// Response cookies.
		cookies []*http.Cookie
		// Response body.
		body string
	}
)

// Create new application instance for test environment.
func newTestApplication(t *testing.T) *application {
	t.Helper()

	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		formDecoder:    form.NewDecoder(),
		templateCache:  templateCache,
		sessionManager: sessionManager,
		debug:          false,
		repositories: &repositories.Repositories{
			Snippet: &mocks.SnippetRepository{},
			User:    &mocks.UserRepository{},
		},
	}
}

// Create new test http server.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	t.Helper()

	ts := &testServer{httptest.NewTLSServer(h)}

	ts.resetClientCookieJar(t)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return ts
}

// Test http GET request to test server.
func (ts *testServer) get(t *testing.T, urlPath string) testResponse {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL+urlPath, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}

// Reset test server http client cookie jar.
func (ts *testServer) resetClientCookieJar(t *testing.T) {
	t.Helper()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
}
