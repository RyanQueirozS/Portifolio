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
		Title: "Not defined",
	}
)

func renderPage(w http.ResponseWriter, pageTemplate string) {
	tmpl := template.Must(template.ParseFiles("./web/template/resources/" + pageTemplate))
	if tmpl == nil {
		log.Fatal("Template file is nil")
	}

	err := tmpl.ExecuteTemplate(w, pageTemplate, pageData)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func InitHandlers() {
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", handleStartPage)
	http.HandleFunc("/about/", handleAboutPage)
	http.HandleFunc("/copyright/", handleCopyrightPage)
	http.HandleFunc("/posts/", handlePostsPage)
	http.HandleFunc("/devlog/", handlePostsPage)

	http.HandleFunc("/web/static/css/style.css", handleCSSFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStartPage(w http.ResponseWriter, r *http.Request) {
	pageData.Title = "Ryan Queiroz - Welcome"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("startpage/index.md", html.FlagsNone))
	renderPage(w, "layout.html")
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

	renderPage(w, "layout.html")
}

func handlePostsPage(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	var endOfPath string
	endOfPath = pathSegments[len(pathSegments)-1]
	if r.URL.Path[len(r.URL.Path)-1] == '/' {
		endOfPath = pathSegments[len(pathSegments)-2]
	}

	switch endOfPath {
	case "devlog":
		{
			devlogs, err := os.ReadDir("./web/template/devlog")
			if err != nil {
				log.Println("Error reading dir", err)
				return

			}

			pageContent := "<ul>"
			for _, post := range devlogs {
				pageContent += "<li><a href=\"" + post.Name() + "\">" + post.Name()[0:len(post.Name())-3] + "</a></li>\n"
				// magical string operations. it basically creates a ul with a link given the page's name
				// the post.Name()[0:len(post.Name())-3] magic, just removes the '.md'. This is done in the case below as well
			}
			pageContent += "<ul>"

			pageData.Title = "Ryan Queiroz - Devlog"
			pageData.Content = template.HTML(pageContent)
			if len(devlogs) == 0 {
				pageData.Content = "<h1>No pages found</h1>"
			}
		}
	case "posts":
		{
			posts, err := os.ReadDir("./web/template/posts")
			if err != nil {
				log.Println("Error reading dir", err)
				return

			}

			pageContent := "<ul>"
			for _, post := range posts {
				pageContent += "<li><a href=\"" + post.Name() + "\">" + post.Name()[0:len(post.Name())-3] + "</a></li>\n"
			}
			pageContent += "<ul>"

			pageData.Title = "Ryan Queiroz - Posts"
			pageData.Content = template.HTML(pageContent)
			if len(posts) == 0 {
				pageData.Content = "<h1>No pages found</h1>"
			}
		}
	default:
		{
			fullPath := pathSegments[len(pathSegments)-2] + "/" + endOfPath

			pageData.Content = template.HTML(converter.GrabMdFileAsHtml(fullPath, html.FlagsNone))
		}
	}

	renderPage(w, "layout.html")
}

func handleCopyrightPage(w http.ResponseWriter, r *http.Request) {
	pageData.Title = "Ryan Queiroz - Copyright"
	pageData.Content = template.HTML(converter.GrabMdFileAsHtml("copyright/copyright.md", html.HrefTargetBlank))
	renderPage(w, "layout.html")
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
