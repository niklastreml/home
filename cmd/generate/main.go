package main

import (
	"github.com/niklastreml/home/pkg/generator"
	"github.com/niklastreml/home/pkg/pages"
)

func main() {
	g := generator.New()
	g.Layout(pages.Layout)
	g.Add("/", pages.Index())
}
