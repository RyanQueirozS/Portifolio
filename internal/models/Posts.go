package models

import "html/template"

type Post struct {
	Name        template.HTML
	Description template.HTML
}
