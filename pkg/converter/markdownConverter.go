package converter

import (
	"log"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func md2Html(mdFile []byte, flags html.Flags) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdFile)

	htmlFlags := html.CommonFlags | flags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func GrabMdFileAsHtml(filePath string, flags html.Flags) []byte {
	fileContents, err := os.ReadFile("web/template/" + filePath)
	if err != nil {
		log.Fatalln("Couldn't open file: ", err)
	}
	htmlContent := md2Html(fileContents, flags)
	return htmlContent
}
