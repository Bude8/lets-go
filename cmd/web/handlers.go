package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Bude8/lets-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// // Initialize a slice containing the paths to the two files. It's important
	// // to note that the file containing our base template must be the *first*
	// // file in the slice.
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }
	//
	// // Use the template.ParseFiles() function to read the files and store the
	// // templates in a template set. Notice that we use ... to pass the contents
	// // of the files slice as variadic arguments.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Use the ExecuteTemplate() method to write the content of the "base"
	// // template as the response body.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
}

func (app *application) getSnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getSnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Displaying a form for creating a new snippet..."))
}

func (app *application) postSnippetCreate(w http.ResponseWriter, r *http.Request) {
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
