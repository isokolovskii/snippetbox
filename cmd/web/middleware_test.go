package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.isokol.dev/internal/assert"
)

const (
	ok = "OK"
)

func TestCommonHeaders(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(ok))
	})

	commonHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get(cspHeaderKey), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get(reffererPolicyHeaderKey), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get(contentTypeOptionsHeaderKey), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, res.Header.Get(frameOptionsHeaderKey), expectedValue)

	expectedValue = "0"
	assert.Equal(t, res.Header.Get(xssProtectionHeaderKey), expectedValue)

	expectedValue = "Go"
	assert.Equal(t, res.Header.Get(serverHeaderKey), expectedValue)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), ok)
}
