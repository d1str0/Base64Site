package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

const version = "0.0.1"
const templateDir = "templates/*"

// Must be able to compile all template files.
var templates = template.Must(template.ParseGlob(templateDir))

func main() {
	http.HandleFunc("/", HomeHandler)              // Should load login page or forward to /admin/system if logged in
	http.HandleFunc("/resources/", includeHandler) // Loads css/js/etc. straight through.

	http.ListenAndServe(":8080", nil)
}

type Page struct {
	Pt string // Plaintext
	En string // Encoded
}

// TODO: Check for those already logged in and forward to /admin
// Loads the login page.
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

// Render just the login page.
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
