package generator

import (
	"context"
	"fmt"
	"os"

	"path/filepath"

	"github.com/a-h/templ"
)

type Generator struct {
	// maps the path to the templ component
	pages  map[string]templ.Component
	layout Layout
}

type Layout func(children templ.Component) templ.Component

func New() *Generator {
	return &Generator{
		pages: map[string]templ.Component{},
	}
}

// Add adds a page to the generator
func (g *Generator) Add(path string, component templ.Component) {
	g.pages[path] = component
}

// Register a layout which will wrap all pages
// A layout is just a component which takes a component as a child
func (g *Generator) Layout(component Layout) {
	g.layout = component
}

// Generates html from the added pages and writes them to the provided file system
func (g *Generator) Generate(ctx context.Context, out string) error {

	for path, component := range g.pages {
		dest := filepath.Join(out, path)
		f, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileWrite, err)
		}
		if g.layout != nil {
			// wrap the component in the layout, if one is provided
			component = g.layout(component)
		}
		if err := component.Render(ctx, f); err != nil {
			return fmt.Errorf("%w: %v", ErrFileCreate, err)
		}
	}

	return nil
}

var ErrFileWrite = fmt.Errorf("failed to write to file")
var ErrFileCreate = fmt.Errorf("failed to create file")
