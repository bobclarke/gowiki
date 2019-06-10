// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Calling the main function")
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling the viewHandler function")
	title := r.URL.Path[len("/view/"):]
	fmt.Println("Page title is: " + title)
	p, err := loadPage(title)

	if err == nil {
		fmt.Println("Page contents are: " + string(p.Body))
		renderTemplate(w, "view", p)
	} else {
		fmt.Println(err)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling the editHandler function")

	// Get the title from the URL
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling the saveHandler function")
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	fmt.Println("Saving " + p.Title + ".txt with contents " + string(p.Body))
	p.save()
}

func (p *Page) save() error {
	fmt.Println("Calling the save function")
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	fmt.Println("Calling the loadPage function for " + title + ".txt")
	body, err := ioutil.ReadFile(title + ".txt")
	if err != nil {
		return nil, err
	}

	// Update Page struct whilst createing a pointer
	p := &Page{Title: title, Body: body}
	return p, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fmt.Println("Calling the renderTemplate function")
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Println(err)
	} else {
		t.Execute(w, p)
	}
}
