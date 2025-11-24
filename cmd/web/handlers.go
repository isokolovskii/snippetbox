package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.isokol.dev/internal/models"
)

const (
	minID             = 1
	defaultExpiration = 7
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

	snippets, err := app.repositories.Snippet.Latest(request.Context())
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	data := app.newTemplateData(request)
	data.Snippets = snippets

	app.renderTemplate(writer, request, http.StatusOK, "home.tmpl.html", data)
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

	app.renderTemplate(writer, request, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("Create a snippet..."))
	if err != nil {
		app.serverError(writer, request, err)
	}
}

func (app *application) snippetCreatePost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	title := "Dummy snippet"
	content := "Dummy snippet content\nSome other content\nWill be removed later"
	expires := defaultExpiration

	id, err := app.repositories.Snippet.Insert(request.Context(), title, content, expires)
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
