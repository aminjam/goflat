package goflat

import (
	"reflect"
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
			//e.g. map "Name,Age,Job" "|"  => "[John|25|Painter Jane|21|Teacher]"
			"map": func(f, sep string, a interface{}) ([]string, error) {
				fields := strings.Split(f, ",")
				reflectedArray := reflect.ValueOf(a)
				out := make([]string, reflectedArray.Len())
				i := 0
				for i < len(out) {
					v := reflectedArray.Index(i)
					row := make([]string, len(fields))
					for k, field := range fields {
						row[k] = v.FieldByName(field).String()
					}
					out[i] = strings.Join(row, sep)
					i++
				}
				return out, nil
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
