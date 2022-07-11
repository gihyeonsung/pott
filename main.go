package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	title   *string
	content *string
	created *string
	updated *string
}

type file struct {
	name string
}

func read(root string) (*category, error) {
	rootCategory := &category{name: "/"}

	rootAbs, _ := filepath.Abs(root)
	log.Printf("rootAbs=%+v", rootAbs)

	reader :=
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pathRel, _ := filepath.Rel(rootAbs, path)
			// skip reading "." since that is already created by hand
			if pathRel == "." {
				return nil
			}

			names := strings.Split(pathRel, "/")
			dirs := names[:len(names)-1]
			base := names[len(names)-1]

			log.Printf("path=%+v dirs=%+v, base=%+v", path, dirs, base)
			if info.IsDir() {
				rootCategory.insertCategory(dirs, base)
			} else if filepath.Ext(base) == ".md" {
				d := &document{}
				d.name = base
				d.raw, _ = ioutil.ReadFile(path)
				rootCategory.insertDoc(dirs, d)
			} else {
				f := &file{}
				f.name = base
				rootCategory.insertFile(dirs, f)
			}

			return err
		}

	if err := filepath.Walk(rootAbs, reader); err != nil {
		return nil, err
	}
	return rootCategory, nil

}

func build(c *category) error {
	return errors.New("not implemented")
}

func write(c *category) (*category, error) {
	return nil, errors.New("not implemented")
}

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.Printf("flags: *pathContent=%+v *pathBuild=%+v, *pathLayout=%+v, *pathCss=%+v", *pathContent, *pathBuild, *pathLayout, *pathCss)
	c, _ := read(*pathContent)
	log.Printf("%+v", c)
}
