package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/note/view", app.noteView)
	mux.HandleFunc("/note/create", app.noteCreate)
	// Middleware chain containing the standard middleware for the application.
	std := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return std.Then(mux)
}
