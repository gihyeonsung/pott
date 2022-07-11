package main

import (
	"flag"
	"log"
)

func main() {
	pathContent := flag.String("content", "content", "contents location")
	pathBuild := flag.String("output", "build", "output location")
	pathLayout := flag.String("template", "layout.tmpl", "template location")
	pathCss := flag.String("stylesheet", "index.css", "stylesheet location")
	log.Printf("*pathContent=%+v", *pathContent)
	log.Printf("*pathBuild=%+v", *pathBuild)
	log.Printf("*pathLayout=%+v", *pathLayout)
	log.Printf("*pathCss=%+v", *pathCss)
}
