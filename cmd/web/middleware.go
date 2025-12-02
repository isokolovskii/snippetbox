package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

const (
	// Content Security Policy header key.
	cspHeaderKey = "Content-Security-Policy"
	// Content Security Policy common header value.
	cspHeader = "default-src 'self'; style-src 'self'" +
		" fonts.googleapis.com; font-src fonts.gstatic.com"
	// Referrer Policy header key.
	referrerPolicyHeaderKey = "Referrer-Policy"
	// Referrer Policy common header value.
	referrerPolicyHeader = "origin-when-cross-origin"
	// Content Type Options header key.
	contentTypeOptionsHeaderKey = "X-Content-Type-Options"
	// Content Type Options common header value.
	contentTypeOptionsHeader = "nosniff"
	// Frame Options header key.
	frameOptionsHeaderKey = "X-Frame-Options"
	// Frame Options common header value.
	frameOptionsHeader = "deny"
	// XSS Protection header key.
	xssProtectionHeaderKey = "X-XSS-Protection"
	// XSS Protection common header value.
	xssProtectionHeader = "0"
	// Server header key.
	serverHeaderKey = "Server"
	// Servercommon header value.
	serverHeader = "Go"
)

// ErrServerUnexpected - error for unexpected things happening.
var ErrServerUnexpected = errors.New("unexpected server error")

// Server common headers middleware.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(cspHeaderKey,
			cspHeader)

		writer.Header().Set(referrerPolicyHeaderKey, referrerPolicyHeader)
		writer.Header().Set(contentTypeOptionsHeaderKey, contentTypeOptionsHeader)
		writer.Header().Set(frameOptionsHeaderKey, frameOptionsHeader)
		writer.Header().Set(xssProtectionHeaderKey, xssProtectionHeader)

		writer.Header().Set(serverHeaderKey, serverHeader)

		next.ServeHTTP(writer, request)
	})
}

// Log server requests middleware.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			ip     = request.RemoteAddr
			proto  = request.Proto
			method = request.Method
			uri    = request.URL.RequestURI()
		)

		app.logger.InfoContext(
			request.Context(),
			"received request",
			slogKeyIP,
			ip,
			slogKeyProto,
			proto,
			slogKeyMethod,
			method,
			slogKeyURI,
			uri,
		)

		next.ServeHTTP(writer, request)
	})
}

// Middleware to recover from panics in handlers.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			pv := recover()

			if pv != nil {
				writer.Header().Set("Connection", "close")
				app.serverError(writer, request, fmt.Errorf("%w: %v", ErrServerUnexpected, pv))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

// Middleware to require authentication for requests.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !app.isAuthenticated(request) {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)

			return
		}

		writer.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(writer, request)
	})
}

// Middleware to prevent cross-site attacks.
func preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// Middleware for checking authentication status and adding it into request context.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := app.sessionManager.GetInt(request.Context(), sessionAuthenticatedUserField)
		if id == 0 {
			next.ServeHTTP(writer, request)

			return
		}

		exists, err := app.repositories.User.Exists(request.Context(), id)
		if err != nil {
			app.serverError(writer, request, err)

			return
		}

		if exists {
			ctx := context.WithValue(request.Context(), isAuthenticatedContextKey, true)
			request = request.WithContext(ctx)
		}

		next.ServeHTTP(writer, request)
	})
}
