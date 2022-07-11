package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
	created string
	updated string
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
				d.title = "undefined"
				d.content = "undefined"
				d.created = "undefined"
				d.updated = "undefined"
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

func write(c *category, out string) error {
	log.Printf("write: dumping c=%+v out=%+v", c, out)
	dir := filepath.Join(out, c.name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	for _, doc := range c.docs {
		path := filepath.Join(dir, doc.name+".html")
		if err := os.WriteFile(path, []byte(doc.content), 0777); err != nil {
			return err
		}
	}

	for _, inner := range c.inners {
		if err := write(inner, dir); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.Printf("flags: *pathContent=%+v *pathBuild=%+v, *pathLayout=%+v, *pathCss=%+v", *pathContent, *pathBuild, *pathLayout, *pathCss)
	c, _ := read(*pathContent)
	build(c)
	if err := write(c, *pathBuild); err != nil {
		log.Panicf("main: %+v", err.Error())
	}

	log.Printf("%+v", c.inners[0].docs[0])
}
