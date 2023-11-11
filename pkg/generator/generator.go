package generator

import (
	"context"
	"fmt"
	"os"
	"strings"

	"path/filepath"

	"github.com/a-h/templ"
)

type Generator struct {
	// maps the path to the templ component
	pages  map[string]templ.Component
	layout Layout
	// maps the url to the path on the file system
	staticFiles map[string]string
}

type Layout func(children templ.Component) templ.Component

func New() *Generator {
	return &Generator{
		pages:       map[string]templ.Component{},
		staticFiles: map[string]string{},
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

	// render pages
	for path, component := range g.pages {
		dest := filepath.Join(out, path)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("%w: %v", ErrFolderCreate, err)
		}

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

	// copy static files
	for url, path := range g.staticFiles {
		dest := filepath.Join(out, url)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("%w: %v", ErrFolderCreate, err)
		}
		if err := copyRecursive(path, dest); err != nil {
			return fmt.Errorf("%w: %v", ErrFileWrite, err)
		}
	}

	return nil
}

func (g *Generator) AddStatic(url, path string) {
	g.staticFiles[url] = path

}

var ErrFileWrite = fmt.Errorf("failed to write to file")
var ErrFileRead = fmt.Errorf("failed to read file")
var ErrFileCreate = fmt.Errorf("failed to create file")
var ErrFolderCreate = fmt.Errorf("failed to create parent folders")

func copyRecursive(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read static file: %w: %w", ErrFileRead, err)
		}
		paths := strings.Split(path, string(os.PathSeparator))[1:]
		paths = append([]string{dest}, paths...)

		path = filepath.Join(paths...)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent folders: %w: %w", ErrFolderCreate, err)
		}
		if err := os.WriteFile(path, b, 0644); err != nil {
			return fmt.Errorf("failed to write static file: %w: %w", ErrFileWrite, err)
		}
		return nil
	})

}
