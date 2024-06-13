package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gustavodiasag/notebox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request path exactly matches the root, if it doesn't, send
	// a 404 response to the client.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	notes, err := app.notes.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl.html", &templateData{
		Notes: notes,
	})
}

func (app *application) noteView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	note, err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, http.StatusOK, "view.tmpl.html", &templateData{
		Note: note,
	})
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

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.notes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/note/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}
