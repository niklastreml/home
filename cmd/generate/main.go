package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v56/github"
	"github.com/niklastreml/home/pkg/api"
	"github.com/niklastreml/home/pkg/generator"
	"github.com/niklastreml/home/pkg/pages"
)

func main() {
	g := generator.New()
	g.Layout(pages.Layout)
	g.Add("index.html", pages.Index())

	client := github.NewClient(http.DefaultClient)
	repos, err := api.GetProjects(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("failed to get projects: %v", err))
	}
	g.Add("projects.html", pages.Projects(repos))
	g.AddStatic("/static", "static")
	if err := g.Generate(context.Background(), "tmp/build"); err != nil {
		panic(err)
	}
}
