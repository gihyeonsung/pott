package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func load(root string) (*category, error) {
	rootCategory := &category{name: "/"}

	rootAbs, _ := filepath.Abs(root)
	log.Printf("load: rootAbs=%+v", rootAbs)

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

			log.Printf("load: path=%+v dirs=%+v, base=%+v", path, dirs, base)
			if info.IsDir() {
				rootCategory.insertCategory(dirs, base)
			} else if filepath.Ext(base) == ".md" {
				d := &document{}
				d.name = base
				d.raw, _ = os.ReadFile(path)
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
