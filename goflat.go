package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//Flat struct
type Flat struct {
	BaseDir  string
	mainFile string
	Template string
	Inputs   []struct{ Path, StructName, VarName string }
}

func (f *Flat) cp(file string) (string, error) {
	outFile := filepath.Join(f.BaseDir, filepath.Base(file))
	in, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer in.Close()
	out, err := os.Create(outFile)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return "", err
	}
	return outFile, cerr
}

func (f *Flat) setInputs(files []string) error {
	for k, v := range files {
		file, err := f.cp(v)
		if err != nil {
			return fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
		}
		f.Inputs[k].Path = file
		name := structName(v)
		f.Inputs[k].StructName = name
		f.Inputs[k].VarName = strings.ToLower(name)
	}
	return nil
}

func (f *Flat) mainGo() error {
	outFile := filepath.Join(f.BaseDir, "main.go")
	main, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer main.Close()
	var tmpl = template.Must(template.New("main").Parse(`package main
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)
type Flater interface {
	Flat() (interface{}, error)
}
func checkError(err error, detail string) {
	if err != nil {
		fmt.Printf("Fatal error %s: %s ", detail, err.Error())
		os.Exit(1)
	}
}
func main() {
	data, err := ioutil.ReadFile("{{.Template}}")
	checkError(err, "reading template file")
	tmpl, err := template.New("").Parse(string(data))
	checkError(err, "parsing template file")
	var result struct {
{{if gt (len .Inputs) 0}}
{{range .Inputs}}
		{{.StructName}} {{.StructName}}
{{end}}
{{end}}
	}
{{if gt (len .Inputs) 0}}
{{range .Inputs}}
	{{.VarName}}, err := New{{.StructName}}().Flat()
	checkError(err, "calling New{{.StructName}}().Flat()")
	result.{{.StructName}} = {{.VarName}}
{{end}}
{{end}}
	var output bytes.Buffer
	err = tmpl.Execute(&output, result)
	checkError(err, "executing template output")
	fmt.Println(string(output.Bytes()))
}`))
	if err := tmpl.Execute(main, f); err != nil {
		return err
	}
	f.mainFile = outFile
	return nil
}

//GoRun runs go on the dynamically created main.go with a given stdout and stderr pipe
func (f *Flat) GoRun(outWriter io.Writer, errWriter io.Writer) error {
	out := make([]string, len(f.Inputs)+2)
	out[0], out[1] = "run", f.mainFile
	for k, v := range f.Inputs {
		out[k+2] = v.Path
	}
	cmd := exec.Command("go", out...)
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	return cmd.Run()
}

//NewFlat Initializes a new Flat struct
func NewFlat(baseDir, template string, inputs []string) (*Flat, error) {
	if _, err := os.Stat(baseDir); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}
	if _, err := os.Stat(template); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}
	f := &Flat{
		BaseDir:  baseDir,
		Template: template,
		Inputs:   make([]struct{ Path, StructName, VarName string }, len(inputs)),
	}

	err := f.setInputs(inputs)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "parsing inputs", err.Error())
	}
	err = f.mainGo()
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "creating main.go", err.Error())
	}
	return f, nil
}

func tmpDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

func structName(filename string) string {
	name := filepath.Base(filename)
	name = strings.Title(strings.Split(name, ".")[0])
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)
	return name
}

const (
	//ErrMissingOnDisk Expected error for accessing invalid file or directory
	ErrMissingOnDisk = "(file or directory is missing)"
)