// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("In the main function")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Define a new type - based on a struct
type Page struct {
	Title string
	Body  []byte
}

type HomePage struct {
	Title string
	Body  []string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// This function will get a list if existing pages and display their names
	fmt.Println("In the homeHandler function")
	fileNames, _ := filepath.Glob("*txt")

	// We need to remove the file extension
	var shortNames []string
	for _, filename := range fileNames {
		shortNames = append(shortNames, strings.TrimSuffix(filename, ".txt"))
	}

	h := &HomePage{Title: "home", Body: shortNames}
	renderHomeTemplate(w, "home", h)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In the viewHandler function")
	title := r.URL.Path[len("/view/"):]
	fmt.Println("Page title is: " + title)
	p, err := loadPage(title)

	// If page does not exist redirect to /edit so that we ca create it
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In the editHandler function")
	title := r.URL.Path[len("/edit/"):] // Get the title from the URL
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title} // If loadPage returns an error, we assume the page doesn't exist so we init a new struct and pointer with Title only
	}
	renderTemplate(w, "edit", p) // Send whichever pointer we have (newly initialised or loaded) to renderTemplate
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In the saveHandler function")
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	fmt.Println("Saving " + p.Title + ".txt with contents " + string(p.Body))
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func (p *Page) save() error {
	fmt.Println("In the save function")
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	body, err := ioutil.ReadFile(title + ".txt")
	fmt.Printf("In the loadPage function for %s.txt \n", title)

	fmt.Println(body)

	bytes.ReplaceAll([]byte(body), []byte("^M"), []byte("<br>"))

	fmt.Println(body)

	if err != nil {
		return nil, err
	}
	p := &Page{Title: title, Body: body} // Update Page struct whilst creating a pointer
	return p, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fmt.Println("In the renderTemplate function")
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Println(err)
	} else {
		t.Execute(w, p)
	}
}

func renderHomeTemplate(w http.ResponseWriter, tmpl string, p *HomePage) {
	t, err := template.ParseFiles(tmpl + ".html")

	if err != nil {
		fmt.Println(err)
	} else {
		t.Execute(w, p)
	}

}
