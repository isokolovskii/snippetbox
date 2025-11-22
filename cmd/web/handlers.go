package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

const (
	minID = 1
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	err = tmpl.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		app.serverError(writer, request, err)
	}
}

func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))

	if err != nil || id < minID {
		http.NotFound(writer, request)

		return
	}

	_, err = fmt.Fprintf(writer, "Display a specific snippet with ID %d...", id)
	if err != nil {
		app.serverError(writer, request, err)
	}
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
	writer.WriteHeader(http.StatusCreated)

	_, err := writer.Write([]byte("Save a new snippet..."))
	if err != nil {
		app.serverError(writer, request, err)
	}
}
