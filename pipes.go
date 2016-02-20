package goflat

import (
	"strings"
	"text/template"
)

type Pipes struct {
	Map template.FuncMap
}

func NewPipes() *Pipes {
	return &Pipes{
		Map: template.FuncMap{
			"join": func(sep string, a []string) (string, error) {
				return strings.Join(a, sep), nil
			},
			"toUpper": func(s string) (string, error) {
				return strings.ToUpper(s), nil
			},
			"toLower": func(s string) (string, error) {
				return strings.ToLower(s), nil
			},
		},
	}
}
