package site

import (
	"sync"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

var (
	markdownOnce sync.Once
	markdown     goldmark.Markdown
)

func markdownRenderer() goldmark.Markdown {
	markdownOnce.Do(func() {
		markdown = goldmark.New(
			goldmark.WithExtensions(
				extension.NewLinkify(),
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
			),
		)
	})
	return markdown
}
