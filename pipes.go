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
			"replace": func(old, new, s string) (string, error) {
				//replace all occurrences of a value
				return strings.Replace(s, old, new, -1), nil
			},
			"split": func(sep, s string) ([]string, error) {
				s = strings.TrimSpace(s)
				if s == "" {
					return []string{}, nil
				}
				return strings.Split(s, sep), nil
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
