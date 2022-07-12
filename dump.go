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
	log.WithFields(log.Fields{
		"c.name": c.name,
		"out":    out,
	}).Info("dump")

	if err := os.MkdirAll(out, 0777); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(out, "index.html"), []byte(c.rendered), 0666); err != nil {
		return err
	}

	for _, d := range c.docs {
		path := filepath.Join(out, d.name)
		if err := os.WriteFile(path, []byte(d.rendered), 0666); err != nil {
			return err
		}
	}

	for _, f := range c.files {
		path := filepath.Join(out, f.name)
		if err := os.WriteFile(path, f.raw, 0666); err != nil {
			return err
		}
	}

	for _, i := range c.inners {
		if err := dumpCategory(i, filepath.Join(out, i.name)); err != nil {
			return err
		}
	}

	return nil
}
