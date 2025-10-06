package main

import "net/http"

func (app *application) routes() *http.ServeMux {
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

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView) //static path
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
