package main

import (
	"log"
	"os"
	"path/filepath"
)

func dump(c *category, out string) error {
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
		if err := dump(inner, dir); err != nil {
			return err
		}
	}

	return nil
}
