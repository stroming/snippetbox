package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/stroming/snippetbox/cmd/web/config"
	"github.com/stroming/snippetbox/cmd/web/models"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of
	// and some short help text explaining what the flag controls. The value of
	// flag will be stored in the addr variable at runtime.
	cfg := new(config.FlagsConfig)
	// port := flag.String("port", ":4000", "HTTP network address")
	flag.StringVar(&cfg.Port, "port", ":4000", "HTTP Network address")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This
	// three parameters: the destination to write the logs to (os.Stdout), a st
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the fl
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stde
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	// As I said above, my general recommendation is to log your output to
	// standard streams and redirect the output to a file at runtime. But if you
	// don’t want to do this, you can always open a file in Go and use it as your
	// log destination. As a rough example:
	// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.Handle("/", models.Home(app))
	mux.Handle("/snippet", models.ShowSnippet(app))
	mux.Handle("/snippet/create", models.CreateSnippet(app))

	// mux.HandleFunc("/snippet/create", models.CreateSnippet(app))
	// Create a file server which serves files out of the "./ui/static" directory
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit.
	// infoLog.Printf("Starting server on %s", cfg.port)
	// err := http.ListenAndServe(cfg.port, mux)

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logge
	// the event of any problems.
	srv := &http.Server{
		Addr:     cfg.Port,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", cfg.Port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
