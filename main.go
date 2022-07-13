package main

import (
	log "github.com/sirupsen/logrus"
)

var (
	pathContent = "content"
	pathBuild   = "build"
	pathLayout  = "layout.tmpl"
	pathCss     = "index.css"
	ignores     = []string{`\.git.*`, `\.github.*`}
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
	c, err := load(pathContent)
	if err != nil {
		log.Panicf("%+v", err.Error())
	}

	if err := build(c, pathLayout, pathCss); err != nil {
		log.Panicf("%+v", err.Error())
	}

	if err := dump(c, pathCss, pathBuild); err != nil {
		log.Panicf("%+v", err.Error())
	}
}
