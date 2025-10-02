package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	//use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {

		/* Use the Header().Set() method to add an 'Allow: POST' header to the
		response header map. The first parameter is the header name, and
		 the second parameter is the header value.
		*/
		w.Header().Set("Allow", "POST")
		/*it it's not, use the w.WriteHeader() method to send a 405 status code and
		the w.Write() method to write "Method Not Allowed" response body.
		*/
		w.WriteHeader(405) //can be called once and always should be before you write the response (w.Write())
		w.Write([]byte("Method Not Allowed"))
		return
	}
	w.Write([]byte("create a new snippet..."))
}

/*
path which ends with a trailing slash are "subtree patterns" which are used to catchall
*/
func main() {
	//Use http.NewServeMux() function to initialize a new servermux, then
	//  register handler functions and corresponding URL patterns with the servemux
	mux := http.NewServeMux()
	// mux.HandleFunc("/", home)                    //subtree pattern - catchall which starts with "/"

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView) //static path
	mux.HandleFunc("/snippet/create", snippetCreate)

	/*use the http.listenAndServe() function to start a new web server. we pass in 2 parameters: the TCP network address
	  to listen on (in this case ":4000") and the servemux we just created
	*/

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	//if the http.ListenAndServe() returns an error we use the log.Fatal() function to log the error message and exit
	log.Fatal(err)
}
