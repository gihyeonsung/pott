package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func list(root string, ignores []string) ([]string, error) {
	var paths []string
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if path == root {
			return nil
		}

		path = strings.TrimPrefix(path, root+"/")
		log.Info(path)
		for _, ignore := range ignores {
			match, err := regexp.Match(ignore, []byte(path))
			if err != nil {
				return err
			}

			if match {
				return nil
			}
		}

		paths = append(paths, path)
		return nil
	}

	if err := filepath.Walk(root, walker); err != nil {
		return nil, err
	}

	return paths, nil
}

func mount(c *category, root string, paths []string) error {
	for _, p := range paths {
		pFromRoot := strings.TrimPrefix(p, root+"/")

		names := strings.Split(pFromRoot, "/")
		dirs, base := names[:len(names)-1], names[len(names)-1]

		cur := c
		for _, dir := range dirs {
			cur = cur.getInner(dir)
		}

		raw, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			return err
		}

		if filepath.Ext(base) == ".md" {
			d := &document{
				name:  base,
				raw:   raw,
				title: base,
			}
			cur.docs = append(cur.docs, d)
		} else {
			cur.blobs = append(cur.blobs, &blob{name: base, raw: raw})
		}
	}

	return nil
}

func load(root string) (*category, error) {
	paths, err := list(root, ignores)
	if err != nil {
		return nil, err
	}

	c := &category{name: "/"}
	if err := mount(c, root, paths); err != nil {
		return nil, err
	}

	return c, nil
}
