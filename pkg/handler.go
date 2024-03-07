package pkg

import (
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, file string) {
	err := tmpl.ExecuteTemplate(w, file, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func InitHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", handleStartPage)

	http.HandleFunc("/about/", handleAboutPage)

	http.HandleFunc("/web/static/css/style.css", handleCSSFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStartPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./web/template/startpage/index.html",
		"./web/template/resources/header.html",
	))

	if tmpl == nil {
		log.Fatal("Template file is nil")
	}
	renderTemplate(w, tmpl, "index.html")
}

func handleAboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./web/template/about/aboutme.html",
		"./web/template/resources/header.html",
	))

	if tmpl == nil {
		log.Fatal("Template file is nil")
	}
	renderTemplate(w, tmpl, "aboutme.html")
}

func handleCSSFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/css/style.css")
}
