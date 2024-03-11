package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown/html"

	"portifolio/internal/models"
	"portifolio/pkg/converter"
)

type PostPageData struct {
	PageData models.PageData
	Posts    []models.Post
	ShowMore template.HTML
}

var (
	postsToShow    = 10
	postIncrement  = 30
	i              = 0
	ShowMoreButton string

	availablePosts = []models.Post{}

	postPage PostPageData
)

func handlePostsPage(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	endOfPath := pathSegments[len(pathSegments)-1]
	if r.URL.Path[len(r.URL.Path)-1] == '/' {
		endOfPath = pathSegments[len(pathSegments)-2]
	}

	switch endOfPath {
	case "devlog":
		{
			availablePosts = []models.Post{}
			handleAvailablePosts("devlog")
			devlogs, err := os.ReadDir("./web/template/devlog")
			if err != nil {
				log.Println("Error reading dir", err)
				return
			}

			if len(devlogs) == 0 {
				pageData.Content = "<h1>No pages found</h1>"
				break
			}
			pageContent := template.HTML(converter.GrabMdFileAsHtml("devlog/"+devlogs[len(devlogs)-1].Name(), html.FlagsNone))
			pageContent += `<hr> <h1>Checkout my latest Devlogs:</h1>`
			pageData.Title = "Ryan Queiroz - Devlog"
			pageData.Content = template.HTML(pageContent)
			postPage = PostPageData{PageData: pageData, Posts: availablePosts}
			ShowMoreButton = ``
			if devlogs[i] != devlogs[len(devlogs)-1] {
				ShowMoreButton = "<div class=\"show-more\"><a href=\"/devlog/index\">Show more</a></div>"
			}
		}
	case "posts":
		{
			availablePosts = []models.Post{}
			handleAvailablePosts("posts")
			posts, err := os.ReadDir("./web/template/posts")
			if err != nil {
				log.Println("Error reading dir", err)
				return
			}
			if len(posts) == 0 {
				pageData.Content = "<h1>No pages found</h1>"
				break
			}
			pageContent := template.HTML(converter.GrabMdFileAsHtml("posts/"+posts[len(posts)-1].Name(), html.FlagsNone))
			pageContent += `<hr> <h1>Checkout my latest Posts:</h1>`
			pageData.Title = "Ryan Queiroz - Posts"
			pageData.Content = template.HTML(pageContent)
			ShowMoreButton = ``
			if posts[i] != posts[len(posts)-1] {
				ShowMoreButton = "<div class=\"show-more\"><a href=\"/devlog/index\">Show more</a></div>"
			}
		}
	default:
		{
			handleAvailablePosts(pathSegments[len(pathSegments)-2])
			if endOfPath == "index" {
				pageData.Content = ""
				pageData.Title = "Devlog index"
				postPage = PostPageData{PageData: pageData, Posts: availablePosts, ShowMore: ""}
				break
			}
			fullPath := pathSegments[len(pathSegments)-2] + "/" + endOfPath
			pageData.Content = template.HTML(converter.GrabMdFileAsHtml(fullPath, html.HrefTargetBlank))
			renderPage(w)
			return
		}
	}
	postPage = PostPageData{PageData: pageData, Posts: availablePosts, ShowMore: template.HTML(ShowMoreButton)}
	tmpl := template.Must(template.ParseFiles("./web/template/resources/postlayout.html"))
	if tmpl == nil {
		log.Fatal("Template file is nil")
	}

	err := tmpl.ExecuteTemplate(w, "postlayout.html", postPage)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func handleAvailablePosts(path string) {
	devlogs, err := os.ReadDir("./web/template/" + path)
	if err != nil {
		log.Println("Error reading dir", err)
		return
	}
	slices.Reverse(devlogs)
	var newPost []models.Post
	for i = range postsToShow {
		if i == len(devlogs) {
			i--
			break
		}
		postDate := strings.Split(devlogs[i].Name()[:10], "-")
		newPost = append(newPost, models.Post{Name: template.HTML(devlogs[i].Name()), Description: template.HTML("Posted on: " + postDate[0] + " " + postDate[1] + " " + postDate[2])})
	}
	availablePosts = newPost
}
