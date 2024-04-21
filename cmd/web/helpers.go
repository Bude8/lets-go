package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// trace  = string(debug.Stack()) # returns byte slice, convert to string so it's readable
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	// Write the template to the buffer instead of straight to w.
	// If there's an error, call our serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If the template is written to buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to w
	w.WriteHeader(status)

	// w is being passed to a function that takes an io.Writer
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// invalid target destination for Decode() will return error with type *form.InvalidDecoderError
		var InvalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &InvalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}
