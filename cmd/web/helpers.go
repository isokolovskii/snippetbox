package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

const (
	slogKeyMethod = "method"
	slogKeyURI    = "uri"
)

func (app *application) serverError(
	writer http.ResponseWriter,
	request *http.Request,
	err error,
) {
	var (
		method = request.Method
		uri    = request.URL.RequestURI()
	)

	app.logger.ErrorContext(
		request.Context(),
		err.Error(),
		slogKeyMethod,
		method,
		slogKeyURI,
		uri,
	)

	if app.debug {
		trace := string(debug.Stack())
		_, _ = fmt.Printf("%s", trace)
	}

	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
