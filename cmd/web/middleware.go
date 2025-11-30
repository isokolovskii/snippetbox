package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// ErrServerUnexpected - error for unexpected things happening.
var ErrServerUnexpected = errors.New("unexpected server error")

// Server common headers middleware.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		writer.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "deny")
		writer.Header().Set("X-XSS-Protection", "0")

		writer.Header().Set("Server", "Go")

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
