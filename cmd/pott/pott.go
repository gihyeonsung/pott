package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gihyeonsung/pott/internal/parser"
	"github.com/gihyeonsung/pott/internal/renderer"
)

const (
	inputDir  = "input"
	outputDir = "output"
)

func build(filename string) error {
	inputPath := path.Join(inputDir, filename)
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return errors.New("could not read the document: " + err.Error())
	}

	doc, err := parser.Parse(inputFile)
	if err != nil {
		return errors.New("could not parse the document: " + err.Error())
	}

	outputPath := path.Join(outputDir, filename+".html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.New("could not create an output file: " + err.Error())
	}

	fmt.Fprintf(outputFile, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="index.css" type="text/css">
</head>
<body>`)
	renderer.Render(outputFile, doc)
	fmt.Fprintf(outputFile, `</body>
</html>`)
	return nil
}

func main() {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read the input directory: %s", err.Error())
		os.Exit(1)
	}

	for _, file := range files {
		build(file.Name())
	}
}
