package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.isokol.dev/internal/models"
)

type (
	snippetCreateForm struct {
		FieldErrors map[string]string
		Title       string
		Content     string
		Expires     int
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
)

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

func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = snippetCreateForm{
		Expires: ExpiresInYear,
	}

	app.renderTemplate(writer, request, http.StatusOK, createTemplateName, data)
}

func (app *application) snippetCreatePost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	form := snippetCreateForm{
		Title:       request.PostForm.Get("title"),
		Content:     request.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	validateSnippetCreateForm(&form)

	if len(form.FieldErrors) > 0 {
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

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
