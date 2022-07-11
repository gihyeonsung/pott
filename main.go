package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func readContentHandler(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	log.Printf("path=%+v, info.IsDir()=%+v", path, info.IsDir())
	return err
}

func readContents(root string) (string, error) {
	rootAbs, _ := filepath.Abs(root)
	log.Printf("rootAbs=%+v", rootAbs)
	filepath.Walk(rootAbs, readContentHandler)
	return "", nil
}

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.Printf("*pathContent=%+v", *pathContent)
	log.Printf("*pathBuild=%+v", *pathBuild)
	log.Printf("*pathLayout=%+v", *pathLayout)
	log.Printf("*pathCss=%+v", *pathCss)

	readContents(*pathContent)
}
