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
	pages map[string]templ.Component
}

func New() *Generator {
	return &Generator{
		pages: map[string]templ.Component{},
	}
}

func (g *Generator) Add(path string, component templ.Component) {
	g.pages[path] = component
}

// Generates html from the added pages and writes them to the provided file system
func (g *Generator) Generate(ctx context.Context, out string) error {

	for path, component := range g.pages {
		dest := filepath.Join(out, path)
		f, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFileWrite, err)
		}
		if err := component.Render(ctx, f); err != nil {
			return fmt.Errorf("%w: %v", ErrFileCreate, err)
		}
	}

	return nil
}

var ErrFileWrite = fmt.Errorf("failed to write to file")
var ErrFileCreate = fmt.Errorf("failed to create file")
