package main

import (
	"log"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func buildDocument(d *document, layout, css string) error {
	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	ctx := parser.NewContext()
	ast := markdown.Parser().Parse(text.NewReader(d.raw), parser.WithContext(ctx))
	metadata := meta.Get(ctx)

	d.title = "no title"

	var content strings.Builder
	markdown.Renderer().Render(&content, d.raw, ast)
	d.content = content.String()

	log.Printf("metamadata=%+v", metadata)
	if created, ok := metadata["created"]; ok {
		if created, ok := created.(string); ok {
			d.created = created
		}
	}

	if updated, ok := metadata["updated"]; ok {
		if updated, ok := updated.(string); ok {
			d.created = updated
		}
	}

	return nil
}

func build(c *category, layout, css string) error {
	for _, i := range c.inners {
		if err := build(i, layout, css); err != nil {
			return err
		}
	}

	for _, d := range c.docs {
		if err := buildDocument(d, layout, css); err != nil {
			return err
		}
	}

	return nil
}
