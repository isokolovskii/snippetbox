package main

import (
	"net/http"
)

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
