package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

func (*application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
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

	buf := new(bytes.Buffer)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	writer.WriteHeader(status)

	_, err = buf.WriteTo(writer)
	if err != nil {
		app.serverError(writer, request, err)
	}
}

func (*application) newTemplateData(_ *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) decodePostForm(request *http.Request, destination any) error {
	err := request.ParseForm()
	if err != nil {
		return fmt.Errorf("error parsing form: %w", err)
	}

	err = app.formDecoder.Decode(destination, request.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return fmt.Errorf("unable to decode form: %w", err)
	}

	return nil
}
