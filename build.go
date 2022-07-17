package main

import (
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	log "github.com/sirupsen/logrus"
)

type layoutParams struct {
	Title    string
	Date     string
	Children []string
	Body     template.HTML
}

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
	log.WithField("c.name", c.name).Info("build")

	// index
	children := []string{".", ".."}
	for _, i := range c.inners {
		children = append(children, i.name+"/")
	}
	for _, d := range c.docs {
		children = append(children, d.name)
	}
	for _, b := range c.blobs {
		children = append(children, b.name)
	}

	var rendered strings.Builder
	err := tmpl.Execute(&rendered, &layoutParams{
		Title:    c.name,
		Children: children,
	})
	if err != nil {
		return err
	}
	c.rendered = rendered.String()

	// docs
	for _, d := range c.docs {
		if err := buildDocument(d, tmpl, cssPath); err != nil {
			return err
		}
	}

	// inner cats
	for _, i := range c.inners {
		if err := buildCategory(i, tmpl, cssPath); err != nil {
			return err
		}
	}

	return nil
}

func buildDocument(d *document, tmpl *template.Template, css string) error {
	log.WithField("d.name", d.name).Info("build")

	markdown := goldmark.New(goldmark.WithExtensions(meta.Meta))
	ctx := parser.NewContext()
	parsed := markdown.Parser().Parse(text.NewReader(d.raw), parser.WithContext(ctx))
	metadata := meta.Get(ctx)

	if date, ok := metadata["date"]; ok {
		if date, ok := date.(string); ok {
			d.date = date
		}
	}
	if d.date == "" {
		log.WithField("d.name", d.name).Warn("no date")
	}

	var body strings.Builder
	markdown.Renderer().Render(&body, d.raw, parsed)

	var rendered strings.Builder
	err := tmpl.Execute(&rendered, &layoutParams{
		Title: d.title,
		Date:  d.date,
		Body:  template.HTML(body.String()),
	})
	if err != nil {
		return err
	}
	d.rendered = rendered.String()

	return nil
}
