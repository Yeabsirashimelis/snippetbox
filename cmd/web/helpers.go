package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

//  STACK TRACE - A stack trace is basically a report of the active function calls in your program at the point when an error or panic happened.
//                  being able to see the execution path of the application via the stack trace can be helpful when you are trying to debug errors.

// http.statusText - automatically generate human friendly text representation of a give http statuscode

// the serverError helper writes an error message and stack trace to the errorLog, then sends a generic 500
//  Internal Server Error response to the user

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) //Sprintf is like Printf, BUT instead of printing to the console, it returns the formatted string
	// app.errorLog.Println(trace)
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description to the user. we'll use this later in book to send response
//  like 400 "Bad Request" when there is a problem with the request that the user sent

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// for consistency we'l also implement a notfound helper. this is actually a convenience wrapper around clientError which sends
//
//	a 404 not found response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
