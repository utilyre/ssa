package view

import (
	"html/template"
	"log/slog"
	"os"
)

func New(l *slog.Logger) *template.Template {
	tmpl, err := template.ParseGlob("views/*.html")
	if err != nil {
		l.Error("failed to parse views", "error", err)
		os.Exit(1)
	}

	return tmpl
}
