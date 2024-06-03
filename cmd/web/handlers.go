package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request path exactly matches the root, if it doesn't, send
	// a 404 response to the client.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// The file containing the base template must be the first one in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Writes the template content as the response body. The last parameter represents
	// any dynamic data that can be additionally passed in.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) noteView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific note with ID %d", id)
}

func (app *application) noteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Must be done before the calls to `w.WriteHeadere()` and `w.Write()` or else
		// there will be no effect on the headers that a user receives.
		w.Header().Set("Allow", http.MethodPost)
		// Implicitly calls `w.WriteHeader()` and `w.Write()` with the respective
		// response status and message.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create new note"))
}
