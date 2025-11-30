package main

import (
	"net/http"

	"github.com/justinas/alice"

	"snippetbox.isokol.dev/ui"
)

const (
	// Route for static files.
	staticRoute = "/static/"
	// Route for home page.
	homeRoute = "/"
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
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET "+staticRoute, http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	mux.Handle("GET "+homeRoute+"{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET "+snippetViewRoute+"/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET "+userSignupRoute, dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST "+userSignupRoute, dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET "+userLoginRoute, dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST "+userLoginRoute, dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET "+snippetCreateRoute, protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST "+snippetCreateRoute, protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST "+userLogoutRoute, protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
