package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)


func home(w http.ResponseWriter, r *http.Request) {
	// If we don't want / to act as catch all
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost) // needs to be called before Write and WriteHeader
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // calls Write and WriteHeader for us
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
