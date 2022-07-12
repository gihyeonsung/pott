package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
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
	log.WithField("c.name", c.name).Info("dump")

	dir := filepath.Join(out, c.name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(c.rendered), 0666); err != nil {
		return err
	}

	for _, d := range c.docs {
		path := filepath.Join(dir, d.name)
		if err := os.WriteFile(path, []byte(d.rendered), 0666); err != nil {
			return err
		}
	}

	for _, i := range c.inners {
		if err := dumpCategory(i, dir); err != nil {
			return err
		}
	}

	return nil
}
