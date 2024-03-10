package models

import (
	"html/template"
)

type PageData struct {
	Title   template.HTML
	Content template.HTML
}
