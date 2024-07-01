package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// Assigns a custom handler for 404 responses.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/note/view/:id", app.noteView)
	router.HandlerFunc(http.MethodGet, "/note/create", app.noteCreate)
	router.HandlerFunc(http.MethodPost, "/note/create", app.noteCreatePost)

	// Middleware chain containing the standard middleware for the application.
	std := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return std.Then(router)
}
