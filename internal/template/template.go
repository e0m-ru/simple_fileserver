package template

import (
	"embed"
	"html/template"

	"github.com/e0m-ru/fileserver/internal/config"
)

//go:embed template/*.html
var templates embed.FS

func CollectTemplates() (err error) {
	config.Config.Os.TMPL, err = template.ParseFS(templates, "template/index.html")
	if err != nil {
		return
	}
	return
}
