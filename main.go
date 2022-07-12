package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

type category struct {
	name     string
	inners   []*category
	docs     []*document
	blobs    []*blob
	rendered string
}

func (c *category) getInner(name string) *category {
	for _, inner := range c.inners {
		// if inner with the name found, returns
		if inner.name == name {
			return inner
		}
	}

	// if not found, create and returns new one
	inner := &category{name: name}
	c.inners = append(c.inners, inner)
	return inner
}

type document struct {
	name     string
	raw      []byte
	title    string
	date     string
	rendered string
}

type blob struct {
	name string
	raw  []byte
}

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.WithFields(log.Fields{
		"pathContent": pathContent,
		"pathBuild":   pathBuild,
		"pathLayout":  pathLayout,
		"pathCss":     pathCss,
	}).Info("config flags")

	c, err := load(*pathContent)
	if err != nil {
		log.Panicf("%+v", err.Error())
	}

	if err := build(c, *pathLayout, *pathCss); err != nil {
		log.Panicf("%+v", err.Error())
	}

	if err := dump(c, *pathCss, *pathBuild); err != nil {
		log.Panicf("%+v", err.Error())
	}
}
