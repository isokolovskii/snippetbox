package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

type (
	// File system wrapper.
	neuteredFileSystem struct {
		// File system.
		fs http.FileSystem
	}
)

const (
	// Route for static files.
	staticRoute = "/static/"
	// Route for home page.
	homeRoute = ""
	// Route for snippet view.
	snippetViewRoute = "/snippet/view"
	// Route for snippet creation.
	snippetCreateRoute = "/snippet/create"
	// Route for user signup.
	userSignupRoute = "/user/signup"
	// Route for user login.
	userLoginRoute = "/user/login"
	// Route for user logout.
	userLogoutRoute = "/user/logout"
)

// Server routes configuration.
func (app *application) routes(staticDir string) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(staticDir)})

	mux.Handle("GET "+staticRoute, http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET "+homeRoute+"/{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET "+snippetViewRoute+"/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET "+snippetCreateRoute, dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST "+snippetCreateRoute, dynamic.ThenFunc(app.snippetCreatePost))

	mux.Handle("GET "+userSignupRoute, dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST "+userSignupRoute, dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET "+userLoginRoute, dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST "+userLoginRoute, dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST "+userLogoutRoute, dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}

// Open file system for read.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	file, err := nfs.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		defer file.Close()

		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if !stat.IsDir() {
		return file, nil
	}

	file, err = nfs.openDirectory(path, file)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}

	return file, nil
}

// Open directory on filesystem for read.
func (nfs neuteredFileSystem) openDirectory(path string, file http.File) (http.File, error) {
	index := filepath.Join(path, "index.html")

	indexFile, err := nfs.fs.Open(index)
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to close file: %w", closeErr)
		}

		return nil, fmt.Errorf("failed to open index.html: %w", err)
	}
	defer indexFile.Close()

	return file, nil
}
