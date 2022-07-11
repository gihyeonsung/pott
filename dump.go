package main

import (
	"log"
	"os"
	"path/filepath"
)

func cp(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return err
	}

	bs, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, bs, 0666); err != nil {
		return err
	}

	return nil
}

func dump(c *category, cssPath, out string) error {
	cp(cssPath, filepath.Join(out, filepath.Base(cssPath)))
	return dumpCategory(c, out)
}

func dumpCategory(c *category, out string) error {
	log.Printf("write: dumping c=%+v out=%+v", c, out)
	dir := filepath.Join(out, c.name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	for _, doc := range c.docs {
		path := filepath.Join(dir, doc.name+".html")
		if err := os.WriteFile(path, []byte(doc.content), 0666); err != nil {
			return err
		}
	}

	for _, inner := range c.inners {
		if err := dumpCategory(inner, dir); err != nil {
			return err
		}
	}

	return nil
}