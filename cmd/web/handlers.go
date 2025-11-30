package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.isokol.dev/internal/models"
	"snippetbox.isokol.dev/internal/validator"
)

type (
	snippetCreateForm struct {
		validator.Validator `form:"-"`

		Title   string `form:"title"`
		Content string `form:"content"`
		Expires int    `form:"expires"`
	}
)

const (
	minID                = 1
	TitleLengthLimit     = 100
	ExpiresInDay         = 1
	ExpiresInWeek        = 7
	ExpiresInYear        = 365
	ValidationErrorBlank = "This field cannot be blank"
	homeTemplateName     = "home.tmpl.html"
	viewTemplateName     = "view.tmpl.html"
	createTemplateName   = "create.tmpl.html"
	fieldTitle           = "title"
	fieldContent         = "content"
	fieldExpires         = "expires"
)

// Handler for home page.
func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := app.repositories.Snippet.Latest(request.Context())
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	data := app.newTemplateData(request)
	data.Snippets = snippets

	app.renderTemplate(writer, request, http.StatusOK, homeTemplateName, data)
}

// Handler for snippet view page.
func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))

	if err != nil || id < minID {
		http.NotFound(writer, request)

		return
	}

	snippet, err := app.repositories.Snippet.Get(request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(writer, request)
		} else {
			app.serverError(writer, request, err)
		}

		return
	}

	data := app.newTemplateData(request)
	data.Snippet = &snippet

	app.renderTemplate(writer, request, http.StatusOK, viewTemplateName, data)
}

// Handler for snippet create page.
func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = snippetCreateForm{
		Expires: ExpiresInYear,
	}

	app.renderTemplate(writer, request, http.StatusOK, createTemplateName, data)
}

// Handler for snippet creation request.
func (app *application) snippetCreatePost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var form snippetCreateForm

	err := app.decodePostForm(request, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	form.validate()

	if !form.Valid() {
		data := app.newTemplateData(request)
		data.Form = form
		app.renderTemplate(writer, request, http.StatusUnprocessableEntity, createTemplateName, data)

		return
	}

	id, err := app.repositories.Snippet.Insert(request.Context(), form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	app.sessionManager.Put(request.Context(), sessionFlashField, "Snippet successfully created!")
	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Validate snippet creation form.
func (form *snippetCreateForm) validate() {
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Title,
		fieldTitle,
		ValidationErrorBlank,
	)

	validator.CheckField(
		&form.Validator,
		validator.CreateMaxCharsValidator(TitleLengthLimit),
		form.Title,
		fieldTitle,
		fmt.Sprintf("This field cannot be more than %d characters long", TitleLengthLimit),
	)

	validator.CheckField(
		&form.Validator, validator.CreateNotBlankValidator(), form.Content, fieldContent, ValidationErrorBlank)

	validator.CheckField(
		&form.Validator,
		validator.CreatePermittedValueValidator(ExpiresInDay, ExpiresInWeek, ExpiresInYear),
		form.Expires,
		fieldExpires,
		fmt.Sprintf(
			"This field must be either %d, %d or %d",
			ExpiresInDay,
			ExpiresInWeek,
			ExpiresInYear,
		),
	)
}
