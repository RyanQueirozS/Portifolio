package handler

import (
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
		Content: "",
	}
)

func renderPage(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("./web/template/resources/layout.html"))
	if tmpl == nil {
		log.Fatal("Template file is nil")
	}

	err := tmpl.ExecuteTemplate(w, "layout.html", pageData)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func InitHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	http.HandleFunc("/", handleStartPage)
	http.HandleFunc("/about/", handleAboutPage)
	http.HandleFunc("/copyright/", handleCopyrightPage)
	http.HandleFunc("/posts/", handlePostsPage)
	http.HandleFunc("/devlog/", handlePostsPage)

	http.HandleFunc("/web/static/css/style.css", handleCSSFiles)
	http.HandleFunc("/web/static/js/", handleJsFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStartPage(w http.ResponseWriter, r *http.Request) {
	pageData.Title = "Ryan Queiroz - Welcome"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("startpage/index.md", html.FlagsNone))
	renderPage(w)
}

func handleAboutPage(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	endOfUrl := segments[len(segments)-1]

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
	pageData.Title = "Ryan Queiroz - Copyright"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("copyright/copyright.md", html.HrefTargetBlank))
	renderPage(w)
}

func handleCSSFiles(w http.ResponseWriter, r *http.Request) {
	var combinedContent []byte
	cssFiles, err := os.ReadDir("web/static/css/")
	if err != nil {
		log.Fatal(err)
	}

	for _, cssFile := range cssFiles {
		fileContent, err := os.ReadFile("web/static/css/" + cssFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		combinedContent = append(combinedContent, fileContent...)
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(combinedContent)
}

func handleJsFiles(w http.ResponseWriter, r *http.Request) {
	var combinedContent []byte
	jsFiles, err := os.ReadDir("web/static/js/")
	if err != nil {
		log.Fatal(err)
	}

	for _, cssFile := range jsFiles {
		fileContent, err := os.ReadFile("web/static/js/" + cssFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		combinedContent = append(combinedContent, fileContent...)
	}

	w.Write(combinedContent)

}
