package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown/html"

	"portifolio/internal/models"
	"portifolio/pkg/converter"
)

var (
	pageData models.PageData = models.PageData{
		Title:   "Not defined",
		Heading: template.HTML(grabHeader()),
		Footer:  template.HTML(grabFooter()),
	}
)

func renderPage(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("./web/template/layout.html"))
	if tmpl == nil {
		log.Fatal("Template file is nil")
	}

	err := tmpl.ExecuteTemplate(w, "layout.html", pageData)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func InitHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", handleStartPage)
	http.HandleFunc("/about/", handleAboutPage)
	http.HandleFunc("/copyright/", handleCopyrightPage)

	http.HandleFunc("/web/static/css/style.css", handleCSSFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStartPage(w http.ResponseWriter, r *http.Request) {
	pageData.Heading = template.HTML(grabHeader())
	pageData.Title = "Ryan Queiroz - Welcome"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("startpage/index.md", html.FlagsNone))
	pageData.Footer = template.HTML(grabFooter())
	renderPage(w)
}

func handleAboutPage(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	endOfUrl := segments[len(segments)-1]

	pageData.Heading = template.HTML(grabHeader())
	pageData.Footer = template.HTML(grabFooter())

	switch endOfUrl {
	default:
		{
			pageData.Title = "Ryan Queiroz - Index"
			pageData.Content = template.HTML(converter.GrabMdFileAsHtml("about/index.md", html.FlagsNone))
		}
	case "me":
		{
			pageData.Title = "Ryan Queiroz - About me"
			pageData.Content = template.HTML(converter.GrabMdFileAsHtml("about/me.md", html.FlagsNone))
		}
	case "my-stack":
		{
			pageData.Title = "Ryan Queiroz - About my stack"
			pageData.Content = template.HTML(converter.GrabMdFileAsHtml("about/my-stack.md", html.HrefTargetBlank))
		}
	}

	renderPage(w)
}

func handleCopyrightPage(w http.ResponseWriter, r *http.Request) {
	pageData.Heading = template.HTML(grabHeader())
	pageData.Title = "Ryan Queiroz - Copyright"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("copyright/copyright.md", html.HrefTargetBlank))
	pageData.Footer = template.HTML(grabFooter())
	renderPage(w)
}

func handleCSSFiles(w http.ResponseWriter, r *http.Request) {
	styleContent, err := os.ReadFile("web/static/css/style.css")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading style.css: %s", err), http.StatusInternalServerError)
		return
	}

	siteHeaderContent, err := os.ReadFile("web/static/css/header.css")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading site-header.css: %s", err), http.StatusInternalServerError)
		return
	}

	pageContent, err := os.ReadFile("web/static/css/page-content.css")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading page-content.css: %s", err), http.StatusInternalServerError)
		return
	}

	siteFooter, err := os.ReadFile("web/static/css/footer.css")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading page-content.css: %s", err), http.StatusInternalServerError)
		return
	}

	combinedContent := append(styleContent, siteHeaderContent...)
	combinedContent = append(combinedContent, pageContent...)
	combinedContent = append(combinedContent, siteFooter...)

	w.Header().Set("Content-Type", "text/css")

	w.Write(combinedContent)
}

func grabHeader() string {
	fileContents, err := os.ReadFile("web/template/resources/header.html")
	if err != nil {
		log.Fatalln("Couldn't open file: ", err)
	}
	return string(fileContents)
}

func grabFooter() string {
	fileContents, err := os.ReadFile("web/template/resources/footer.html")
	if err != nil {
		log.Fatalln("Couldn't open file: ", err)
	}
	return string(fileContents)
}
