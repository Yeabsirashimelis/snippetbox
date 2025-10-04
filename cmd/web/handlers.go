package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// define a home handler function which writes a byte slice containing "Hello from the snippetbox" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	/*
		   for instance, in the application we're building  we want the home page to be displayed if and only if - the request URL path
		     exactly matches "/". otherwise, we want the user to receive a 404 page not found response

			 b/c we don't want the "/" to be a catch-all
	*/
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
		log.Println(err.Error()) // log the error as a string
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	/*
	   use the Execute template() method to write the content of the "base"
	   //template as the response body
	*/
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	/* Extract the value of the id parameter from the query string and try to
	   convert it to an integer using strconv.Atoi() function. if it can't be converted to an integer,
	   or the value less than 1, return a 404 page not found response
	*/

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to http.ResponseWritter
	// fmt.Fprintf - you can write using this literally anywhere - Writing to console (same as fmt.Printf), writting to an http response, writting to a file
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("create a new snippet..."))
}
