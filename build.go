package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func build(c *category, layoutPath, cssPath string) error {
	layout, err := os.ReadFile(layoutPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("layout").Parse(string(layout))
	if err != nil {
		return err
	}

	return buildCategory(c, tmpl, cssPath)
}

func buildCategory(c *category, tmpl *template.Template, cssPath string) error {
	for _, i := range c.inners {
		if err := buildCategory(i, tmpl, cssPath); err != nil {
			return err
		}
	}

	for _, d := range c.docs {
		if err := buildDocument(d, tmpl, cssPath); err != nil {
			return err
		}
	}

	return nil
}

func buildDocument(d *document, tmpl *template.Template, css string) error {
	log.Printf("buildDocument: building d.name=%+v", d.name)

	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	ctx := parser.NewContext()
	parsed := markdown.Parser().Parse(text.NewReader(d.raw), parser.WithContext(ctx))
	metadata := meta.Get(ctx)

	if parsed.HasChildren() {
		node := parsed.FirstChild()
		if node.Kind() == ast.KindHeading {
			d.title = string(node.Text(d.raw))
		}
	}

	if date, ok := metadata["date"]; ok {
		if date, ok := date.(string); ok {
			d.date = date
		}
	}

	var body strings.Builder
	markdown.Renderer().Render(&body, d.raw, parsed)

	var content strings.Builder
	tmpl.Execute(&content, &struct {
		Title string
		Date  string
		Body  template.HTML
	}{
		Title: d.title,
		Date:  d.date,
		Body:  template.HTML(body.String()),
	})
	d.content = content.String()

	return nil
}
