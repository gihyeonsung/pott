package main

import (
	"log"
	"os"
	"path/filepath"
)

func dump(c *category, cssPath, out string) error {
	css, err := os.ReadFile(cssPath)
	if err != nil {
		return err
	}

	cssPathOut := filepath.Join(out, filepath.Base(cssPath))
	if err := os.WriteFile(cssPathOut, []byte(css), 0777); err != nil {
		return err
	}

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
		if err := os.WriteFile(path, []byte(doc.content), 0777); err != nil {
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
