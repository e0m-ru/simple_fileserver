package main

import (
	"embed"
	"html/template"
)

//go:embed template/*.html
var templates embed.FS

func collectTemplates(cfg *Config) (err error) {
	cfg.os.TMPL, err = template.ParseFS(templates, "template/index.html")
	if err != nil {
		return
	}
	return
}
