package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

const version = "0.0.1"
const templateDir = "templates/*"

// Must be able to compile all template files.
var templates = template.Must(template.ParseGlob(templateDir))

func main() {
	http.HandleFunc("/", HomeHandler)              // Load main page
	http.HandleFunc("/resources/", includeHandler) // Loads css/js/etc. straight through.

	srv := &http.Server{
		Addr:         ":443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}

type Page struct {
	Pt string // Plaintext
	En string // Encoded
}

// Handle all requests for home page, including POST.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%q\n", r)
	if r.Method != "POST" {
		renderIndex(w, &Page{})
	} else {
		plaintext := r.PostFormValue("plaintext")
		if plaintext == "" {
			encoded := r.PostFormValue("encoded")
			if encoded == "" {
				renderIndex(w, &Page{})
			} else {
				page := &Page{}
				temp, err := base64.StdEncoding.DecodeString(encoded)
				page.Pt = string(temp)
				if err != nil {
					page.Pt = err.Error()
				}
				renderIndex(w, page)
			}
		} else {
			page := &Page{}
			page.En = base64.StdEncoding.EncodeToString([]byte(plaintext))
			renderIndex(w, page)
		}
	}
}

// Render just the home page.
func renderIndex(w http.ResponseWriter, p *Page) {
	err := templates.ExecuteTemplate(w, "main", p)
	if err != nil {
		panic(err.Error())
	}
}

// For resource files like js, images, etc.
// Just a straight through file server.
func includeHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[1:]
	http.ServeFile(w, r, filename)
}
