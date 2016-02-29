package main

import (
	"strings"
	"text/template"
)

func CustomPipes() template.FuncMap {
	return template.FuncMap{
		//This will override the default "replace" pipe.
		"replace": func(old, new, s string) (string, error) {
			//replace only the first occurrence of a value
			return strings.Replace(s, old, new, 1), nil
		},
		//This will extend the list of helper functions
		"sanitize": func(a string) (string, error) {
			if a == "SECRET" {
				return "TERCES", nil
			}
			return a, nil
		},
	}
}
