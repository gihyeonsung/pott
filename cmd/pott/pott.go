package main

import (
	"fmt"
	"os"

	"github.com/gihyeonsung/pott/internal/parser"
	"github.com/gihyeonsung/pott/internal/renderer"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read the document: %s", err.Error())
		os.Exit(1)
	}

	doc, err := parser.Parse(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse the document: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="index.css" type="text/css">
</head>
<body>`)
	renderer.Render(os.Stdout, doc)
	fmt.Printf(`</body>
</html>`)
}
