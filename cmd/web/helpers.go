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

const (
	// Field saved in session for showing flash messages.
	sessionFlashField = "flash"
	// Field saved in session for user id.
	sessionAuthenticatedUserField = "authenticatedUserID"
)

// ErrTemplateNotFound - error returned if required template not found.
var ErrTemplateNotFound = errors.New("template not found")

// Helper for returning server error to user.
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

// Helper for returning client error to user.
func (*application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

// Helper function for rendering templates.
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

// Helper function for creating template data.
func (app *application) newTemplateData(request *http.Request) *templateData {
	flash := app.sessionManager.PopString(request.Context(), sessionFlashField)
	isAuthenticated := app.isAuthenticated(request)

	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           flash,
		IsAuthenticated: isAuthenticated,
	}
}

// Helper function for decoding forms in post requests.
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

// Checks if user is authenticated.
func (app *application) isAuthenticated(request *http.Request) bool {
	return app.sessionManager.Exists(request.Context(), sessionAuthenticatedUserField)
}
