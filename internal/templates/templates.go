package templates

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"

	"go.uber.org/fx"
)

type Templates struct {
	Pages *template.Template
}

func New(lc fx.Lifecycle, l *slog.Logger) Templates {
	pages := template.New("layout.html")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if _, err := pages.
				Funcs(newFuncMap(pages)).
				ParseGlob("pages/*.html"); err != nil {
				return err
			}

			return nil
		},
	})

	return Templates{
		Pages: pages,
	}
}

func newFuncMap(tmpl *template.Template) template.FuncMap {
	return template.FuncMap{
		"page": func(name string, data any) (template.HTML, error) {
			buf := new(bytes.Buffer)
			if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	}
}
