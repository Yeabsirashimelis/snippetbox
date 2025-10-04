package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

/*
path which ends with a trailing slash are "subtree patterns" which are used to catchall
*/
func main() {

	/*
			 Define a new command-line flag with the name "addr", a default value of ":4000" and some short help
			  text explaining what the flag controls. the value of the flag will be stored in the addr variable at runtime.

		  flag.string() - this has a benefit of converting whatever value the user provides at runtime to a string type.
		                   if the value can't be converted to a string then the application will log an error and exit

		  go has a range of other functions including flag.Int(), flag.Bool(),... - they automatically convert the command-line flg to the appropriate type.
	*/
	addr := flag.String("addr", ":4000", "HTTP network address")

	/*
		   Importantly, we use the flag.parse() function to parse the commad-line flag. this reads in the command-line flag value
		    and assigns it to the addr variable. you need to call this **before** you use the addr variable otherwise it will always
			contain the default value ":4000". if any errors are encountered during parsing the aplication will be terminated
	*/

	flag.Parse()

	/*
		    use the log.new() to create a logger for writting information messages. this takes three parameters: the destination to write logs to (os.Stdout), a string
			prefix for message (INFO followed by a tab) and flags to indicate that additional inforamtion to include (local date and time)
			  NOTE THAT: the flags are joined using the bitwise OR operator |.
	*/
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	/*
	   create a logger for writting error messages in the same way, but the stderr as the destination and use the Log.Lshortfile flag to include the relevant
	    file name and line number
	*/
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // Llongfile if you want to include the full-path

	//Use http.NewServeMux() function to initialize a new servermux, then
	//  register handler functions and corresponding URL patterns with the servemux
	mux := http.NewServeMux()
	// mux.HandleFunc("/", home)                    //subtree pattern - catchall which starts with "/"

	// create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project directory root
	fileServer := http.FileServer(http.Dir("./ui/static/")) // "hey, here is the directory that contains the fiels you can serve and create http handler that can serve static files from that directory"

	// use the mux.handle() function to register the file server as the handler for all URL paths that start with "/static/".
	//  for matching paths, we strip the "/static" prefix before the request reaches the file server
	/*
		     /static/ in the URL
		    Any request starting with /static/ will be handled by this handler.
		    http.StripPrefix("/static", fileServer)
			Removes /static from the URL before passing it to fileServer.
			For example:
			URL: /static/css/style.css
			Stripped URL: /css/style.css
			File server looks for: ./ui/static/css/style.css on disk.
	*/
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView) //static path
	mux.HandleFunc("/snippet/create", snippetCreate)

	/*use the http.listenAndServe() function to start a new web server. we pass in 2 parameters: the TCP network address
	  to listen on (in this case ":4000") and the servemux we just created
	*/

	/*
	   the value returned from the flag.String() function is a pointer to the flag value, not the value itself. so we need
	     to dereference the pointer
	*/

	// log.Printf("Starting server on %s", *addr)
	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)

	//if the http.ListenAndServe() returns an error we use the log.Fatal() function to log the error message and exit
	// log.Fatal(err)

	errorLog.Fatal(err)
}
