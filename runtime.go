package goflat

const (
	MainGotempl = `package main
import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "text/template"
    )
func checkError(err error, detail string) {
  if err != nil {
    fmt.Printf("Fatal error %s: %s ", detail, err.Error())
      os.Exit(1)
  }
}
func main() {
  data, err := ioutil.ReadFile("{{.GoTemplate}}")
    checkError(err, "reading template file")
    pipes := NewPipes()
    {{if ne .CustomPipes ""}}
    pipes.Extend(CustomPipes())
    {{end}}
    tmpl, err := template.New("").Funcs(pipes.Map).Parse(string(data))
    checkError(err, "parsing template file")
    var result struct {
      {{if gt (len .GoInputs) 0}}
      {{range .GoInputs}}
      {{.StructName}} {{.StructName}}
      {{end}}
      {{end}}
    }
  {{if gt (len .GoInputs) 0}}
  {{range .GoInputs}}
  result.{{.StructName}} = New{{.StructName}}()
  {{end}}
  {{end}}
  var output bytes.Buffer
    err = tmpl.Execute(&output, result)
    checkError(err, "executing template output")
    fmt.Println(string(output.Bytes()))
}
`
	PipesGo = `package runtime

import (
	"reflect"
	"strings"
	"text/template"
)

type Pipes struct {
	Map template.FuncMap
}

func (p *Pipes) Extend(fm template.FuncMap) {
	for k, v := range fm {
		p.Map[k] = v
	}
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
`
)
