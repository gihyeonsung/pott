package main

import (
	"os"
	"path/filepath"
	"strings"
)

func list(root string) ([]string, error) {
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

		if filepath.Ext(base) == ".md" {
			raw, _ := os.ReadFile(p)
			cur.docs = append(cur.docs, &document{
				name: strings.TrimSuffix(base, ".md"),
				raw:  raw,
			})
		} else {
			cur.files = append(cur.files, &file{name: base})
		}
	}

	return nil
}

func load(root string) (*category, error) {
	c := &category{name: "/"}

	paths, err := list(root)
	if err != nil {
		return nil, err
	}

	if err := mount(c, root, paths); err != nil {
		return nil, err
	}

	return c, nil
}
