package main

import (
	"log"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func build(c *category) error {
	for _, doc := range c.docs {
		markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
		ctx := parser.NewContext()
		ast := markdown.Parser().Parse(text.NewReader(doc.raw), parser.WithContext(ctx))
		metadata := meta.Get(ctx)

		doc.title = "no title"

		var content strings.Builder
		markdown.Renderer().Render(&content, doc.raw, ast)
		doc.content = content.String()

		log.Printf("metamadata=%+v", metadata)
		if created, ok := metadata["created"]; ok {
			if created, ok := created.(string); ok {
				doc.created = created
			}
		}

		if updated, ok := metadata["updated"]; ok {
			if updated, ok := updated.(string); ok {
				doc.created = updated
			}
		}
	}

	for _, inner := range c.inners {
		build(inner)
	}

	return nil
}
