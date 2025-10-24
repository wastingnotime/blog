package site

import (
	"html/template"
)

func Funcs(cfg Config) template.FuncMap {
	return template.FuncMap{
		"href": cfg.Href,
	}
}
