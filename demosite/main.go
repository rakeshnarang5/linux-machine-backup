package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

type pageData struct {
	Title     string
	FirstName string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	//  "/templates/*gohtml" --absolute reference
	//  "templates/gohtml" --relative reference
}

func main() {
	http.HandleFunc("/", idx)
	http.HandleFunc("/about", abot)
	http.HandleFunc("/contact", cntct)
	http.HandleFunc("/apply", aply)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
	// http.ListenAndServe accepts two parameters a string, and a handler like ServeMux
	// when nil is passed DefaultServeMux is used
}

func idx(w http.ResponseWriter, req *http.Request) {
	pd := pageData{
		Title: "Index Page",
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", pd)
	if err != nil {
		log.Println("LOGGED", err)
		http.Error(w, "Internal serverrrrrrrrrrrrrrrr error", http.StatusInternalServerError)
		return
	}
	fmt.Println(req.URL.Path)
	fmt.Println("we got here")
}
func abot(w http.ResponseWriter, req *http.Request) {
	pd := pageData{
		Title: "About Page",
	}
	err := tpl.ExecuteTemplate(w, "about.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("we got here")
}
func cntct(w http.ResponseWriter, req *http.Request) {
	pd := pageData{
		Title: "Contact Page",
	}
	err := tpl.ExecuteTemplate(w, "contact.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("we got here")
}
func aply(w http.ResponseWriter, req *http.Request) {
	pd := pageData{
		Title: "Apply Page",
	}
	var first string
	var last string
	if req.Method == http.MethodPost {
		first = req.FormValue("fname")
		last = req.FormValue("lname")
		pd.FirstName = first + " " + last
	}
	err := tpl.ExecuteTemplate(w, "apply.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("we got here")
}
