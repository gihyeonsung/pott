package main

import (
	"flag"
	"log"
)

type category struct {
	name   string
	outer  *category
	inners []*category
	docs   []*document
	files  []*file
}

func (c *category) getInner(name string) *category {
	for _, inner := range c.inners {
		// if inner with the name found, returns
		if inner.name == name {
			return inner
		}
	}

	// if not found, create and returns new one
	inner := &category{name: name, outer: c}
	c.inners = append(c.inners, inner)
	return inner
}

func (c *category) insertCategory(dirs []string, n string) {
	// log.Printf("insertCategory: dirs=%+v n=%+v", dirs, n)
	cur := c
	for _, dir := range dirs {
		cur = cur.getInner(dir)
	}

	inner := &category{name: n, outer: cur}
	cur.inners = append(cur.inners, inner)
}

func (c *category) insertDoc(dirs []string, d *document) {
	// log.Printf("insertDoc: dirs=%+v d=%+v", dirs, d)
	cur := c
	for _, dir := range dirs {
		cur = cur.getInner(dir)
	}
	cur.docs = append(cur.docs, d)
}

func (c *category) insertFile(dirs []string, f *file) {
	// log.Printf("insertDoc: dirs=%+v f=%+v", dirs, f)
	cur := c
	for _, dir := range dirs {
		cur = cur.getInner(dir)
	}
	cur.files = append(cur.files, f)
}

type document struct {
	name    string
	raw     []byte
	title   string
	content string
	date    string
}

type file struct {
	name string
}

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.Printf("flags: *pathContent=%+v *pathBuild=%+v, *pathLayout=%+v, *pathCss=%+v", *pathContent, *pathBuild, *pathLayout, *pathCss)

	c, err := load(*pathContent)
	if err != nil {
		log.Panicf("main: %+v", err.Error())
	}

	if err := build(c, *pathLayout, *pathCss); err != nil {
		log.Panicf("main: %+v", err.Error())
	}

	if err := dump(c, *pathCss, *pathBuild); err != nil {
		log.Panicf("main: %+v", err.Error())
	}
	log.Printf("all processes complete")
}
