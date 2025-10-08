package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	/*
		Import the models package that we just created. you need to prefix this with  whatever module path you set up back
		 (project setup and creating a module), so that the import statement looks like this :

		 - "{your-module-path}/ internal/ models" - you can find your module path at the top of the go.mod file
	*/
	"snippetbox.yeabsira.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

/*
Define an application struct to hold the application-wide dependencies for the web application. for now we'll only include fields

	for the 2 custom loggers, but we will add more to it as the build processes.
*/
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

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

	// Define a new command-line flag for the MYSQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

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

	/*
	   to keep the main() function tidy, I've put the code for creating a connection pool into the separate openDB() function below
	*/
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	//we also defer a call to db.Close(), so that the connection pool is closed before the main() function exists
	defer db.Close()

	// initialize a new instance of our application struct, containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	/*use the http.listenAndServe() function to start a new web server. we pass in 2 parameters: the TCP network address
	  to listen on (in this case ":4000") and the servemux we just created
	*/

	/*
	   the value returned from the flag.String() function is a pointer to the flag value, not the value itself. so we need
	     to dereference the pointer
	*/

	/*
	   initialize a new http.Server struct. we set the Addr and handler fields so that the server uses the same network address and routes as before,
	     and set the ErrorLog field so that the server now uses the custom errorLog logger in the event of any problems.
	*/

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// log.Printf("Starting server on %s", *addr)
	// infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)

	infoLog.Printf("Starting server on %s", *addr)
	// Call the listenAndServe method on our new http.Server struct
	err = srv.ListenAndServe()

	//if the http.ListenAndServe() returns an error we use the log.Fatal() function to log the error message and exit
	// log.Fatal(err)

	errorLog.Fatal(err)
}

// the openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given dsn(connection string)
func openDB(dsn string) (*sql.DB, error) {
	// the sql.open function initialies a new sql.DB object, which is esentially a pool of database connections
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
