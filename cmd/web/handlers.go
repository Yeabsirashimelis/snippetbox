package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"snippetbox.yeabsira.net/internal/models"
)

// change th signiture of the handler so that it is defined as a method against the *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		   for instance, in the application we're building  we want the home page to be displayed if and only if - the request URL path
		     exactly matches "/". otherwise, we want the user to receive a 404 page not found response

			 b/c we don't want the "/" to be a catch-all
	*/
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w) // use the notFound() helper

		return
	}

	// Initialize a slice containing the paths to the two files. it's important to note that
	// the file containing our base template must be the *first* in the slice.
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	/*
		use the template.ParseFiles() funtion to read the template file into a template set.
		  if there's an error we log the detailed error message and send error response to the client

		  notice that we can pass the slicce of the paths as a variadic parameter?
	*/
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, err) // use the serverError() helper
		return
	}

	/*
	   use the Execute template() method to write the content of the "base"
	   //template as the response body
	*/
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err) // use the serverError() helper

	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	/* Extract the value of the id parameter from the query string and try to
	   convert it to an integer using strconv.Atoi() function. if it can't be converted to an integer,
	   or the value less than 1, return a 404 page not found response
	*/

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w) // use the notFound helper
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)

		} else {
			app.serverError(w, err)
		}
		return
	}

	/*
		// Use the fmt.Fprintf() function to interpolate the id value with our response
		// and write it to http.ResponseWritter
		// fmt.Fprintf - you can write using this literally anywhere - Writing to console (same as fmt.Printf), writting to an http response, writting to a file
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	*/

	fmt.Fprintf(w, "%+v", snippet)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	//use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {

		/* Use the Header().Set() method to add an 'Allow: POST' header to the
		response header map. The first parameter is the header name, and
		 the second parameter is the header value.
		*/
		w.Header().Set("Allow", "POST")

		/*
			it it's not, use the w.WriteHeader() method to send a 405 status code and
			the w.Write() method to write "Method Not Allowed" response body.
		*/
		// w.WriteHeader(405) //can be called once and always should be before you write the response (w.Write())
		// w.Write([]byte("Method Not Allowed"))

		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed) // use the clientError() helper
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
