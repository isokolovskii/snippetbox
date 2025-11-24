package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

const (
	slogKeyMethod = "method"
	slogKeyURI    = "uri"
)

var ErrTemplateNotFound = errors.New("template not found")

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

func (app *application) renderTemplate(
	writer http.ResponseWriter,
	request *http.Request,
	status int,
	page string,
	data *templateData,
) {
	tmpl, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("%w: the template %s does not exist", ErrTemplateNotFound, page)
		app.serverError(writer, request, err)

		return
	}

	writer.WriteHeader(status)

	err := tmpl.ExecuteTemplate(writer, "base", data)
	if err != nil {
		app.serverError(writer, request, err)
	}
}
