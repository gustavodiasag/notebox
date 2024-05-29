package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request path exactly matches the root, if it doesn't, send
	// a 404 response to the client.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Health check"))
}

func noteView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific note with ID %d", id)
}

func noteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Must be done before the calls to `w.WriteHeadere()` and `w.Write()` or else
		// there will be no effect on the headers that a user receives.
		w.Header().Set("Allow", http.MethodPost)
		// Implicitly calls `w.WriteHeader()` and `w.Write()` with the respective
		// response status and message.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create new note"))
}
