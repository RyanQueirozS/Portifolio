package models

import (
	"html/template"
)

type PageData struct {
	Title   template.HTML
	Heading template.HTML
	Content template.HTML
	Footer  template.HTML
}
