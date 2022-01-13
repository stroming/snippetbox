package models

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/stroming/snippetbox/cmd/web/config"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func Home(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Use the template.ParseFiles() function to read the template file into a
		// template set. If there's an error, we log the detailed error message and
		// the http.Error() function to send a generic 500 Internal Server Error
		// response to the user.

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		log.Println("here is the log println")
		// We then use the Execute() method on the template set to write the templa
		// content as the response body. The last parameter to Execute() represents
		// dynamic data that we want to pass in, which for now we'll leave as nil.
		err = ts.Execute(w, nil)
		if err != nil {
			app.ServerError(w, err)
			http.Error(w, "Internal Server Error", 500)
		}
	})
}

func ShowSnippet(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract the value of the id parameter from the query string and try to
		// convert it to an integer using the strconv.Atoi() function. If it can't
		// be converted to an integer, or the value is less than 1, we return a 404
		// not found response.
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.ErrorLog.Println(err.Error())
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	})
}

// Add a createSnippet handler function.
func CreateSnippet(app *config.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Use r.Method to check whether the request is using POST or not.
		// If it's not, use the w.WriteHeader() method to send a 405 status code and
		// the w.Write() method to write a "Method Not Allowed" response body. We
		// then return from the function so that the subsequent code is not executed
		if r.Method != "POST" {

			// Set a new cache-control header. If an existing "Cache-Control" header exists
			// it will be overwritten.
			w.Header().Set("Cache-Control", "public, max-age=31536000")

			// In contrast, the Add() method appends a new "Cache-Control" header and can
			// be called multiple times.
			w.Header().Add("cache-Control", "public123")
			w.Header().Add("Cache-Control", "max-age=31536000123")
			w.Header().Add("Cache-Control2", "max-age=31536000123")

			// Delete all values for the "Cache-Control" header.
			w.Header().Del("Cache-Control2")

			// The Del() method doesn’t remove system-generated headers. To
			// suppress these, you need to access the underlying header map directly
			// and set the value to nil. If you want to suppress the Date header, for
			// example, you need to write:
			w.Header()["Date"] = nil

			// This line will overwrite everything about the previous header
			w.Header().Set("cache-Control", "public, max-age=31536000")

			// Retrieve the first value for the "Cache-Control" header.
			g := w.Header().Get("Cache-Control")
			fmt.Println(g)

			// Use the http.Error() function to send a 405 status code and "Method N
			// Allowed" string as the response body.
			http.Error(w, "Method Not Allowed", 405)
			return
		}

		w.Write([]byte("Create a new snippet..."))
	})
}
