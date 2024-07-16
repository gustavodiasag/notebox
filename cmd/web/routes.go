package main

import (
	"net/http"

	"github.com/gustavodiasag/notebox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// Assigns a custom handler for 404 responses.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/health_check", healthCheck)

	// Middleware chain containing the middleware specific to the dynamic application routes.
	dyn := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dyn.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/note/view/:id", dyn.ThenFunc(app.noteView))
	router.Handler(http.MethodGet, "/user/signup", dyn.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dyn.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dyn.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dyn.ThenFunc(app.userLoginPost))

	// Authenticated-only routes.
	protected := dyn.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/note/create", protected.ThenFunc(app.noteCreate))
	router.Handler(http.MethodPost, "/note/create", protected.ThenFunc(app.noteCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Middleware chain containing the standard middleware for the application.
	std := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return std.Then(router)
}
