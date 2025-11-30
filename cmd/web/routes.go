package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

type (
	neuteredFileSystem struct {
		fs http.FileSystem
	}
)

func (app *application) routes(staticDir string) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(staticDir)})

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}

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
