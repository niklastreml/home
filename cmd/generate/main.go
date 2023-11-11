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
	g.AddStatic("/static", "static")
	if err := g.Generate(context.Background(), "tmp/build"); err != nil {
		panic(err)
	}
}
