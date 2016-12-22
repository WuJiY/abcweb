package app

import (
	"html/template"

	"github.com/nullbio/abcweb/rendering"
	"github.com/unrolled/render"
)

// appHelpers is a map of the template helper functions.
// Assign the template helper funcs in here, example:
// "titleCase": strings.TitleCase
var appHelpers = template.FuncMap{}

// InitRenderer initializes the renderer using the app configuration.
// If you need to use multiple renderers, you can add more Renderer
// variables to your State object and initialize them here.
func (s State) InitRenderer() {
	render := &rendering.Render{
		Render: render.New(render.Options{
			Directory:     s.Config.Templates,
			Layout:        "layout",
			Extensions:    []string{".tmpl", ".html"},
			IsDevelopment: s.Config.RenderRecompile,
			Funcs:         []template.FuncMap{appHelpers},
		}),
	}

	s.Render = render
}
