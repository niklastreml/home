package main

import (
	"context"

	"github.com/niklastreml/home/pkg/generator"
	"github.com/niklastreml/home/pkg/pages"
)

func main() {
	g := generator.New()
	g.Layout(pages.Layout)
	g.Add("index.html", pages.Index())
	if err := g.Generate(context.Background(), "tmp/build"); err != nil {
		panic(err)
	}
}
