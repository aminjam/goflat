package goflat

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
	BaseDir     string
	mainFile    string
	GoTemplate  string
	GoInputs    []struct{ Path, StructName, VarName string }
	GoPipesPath string
}

func (f *Flat) cp(file string) (string, error) {
	base := filepath.Base(file)
	if !strings.HasSuffix(base, ".go") {
		base += ".go"
	}
	outFile := filepath.Join(f.BaseDir, base)
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

func (f *Flat) pipes() error {
	data, err := ioutil.ReadFile("pipes.go")
	if err != nil {
		return err
	}
	content := strings.Replace(string(data), "package goflat", "package main", -1)
	outFile := filepath.Join(f.BaseDir, "pipes.go")
	err = ioutil.WriteFile(outFile, []byte(content), 0666)
	if err != nil {
		return err
	}
	f.GoPipesPath = outFile
	return nil
}

func (f *Flat) setGoInputs(files []string) error {
	for k, v := range files {
		orgFile, structName := splitGoInput(v)
		file, err := f.cp(orgFile)
		if err != nil {
			return fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
		}
		f.GoInputs[k].Path = file
		f.GoInputs[k].StructName = structName
		f.GoInputs[k].VarName = strings.ToLower(structName)
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

	data, err := ioutil.ReadFile("main.gotempl")
	if err != nil {
		return err
	}
	var tmpl = template.Must(template.New("main").Parse(string(data)))
	if err := tmpl.Execute(main, f); err != nil {
		return err
	}
	f.mainFile = outFile
	return nil
}

//GoRun runs go on the dynamically created main.go with a given stdout and stderr pipe
func (f *Flat) GoRun(outWriter io.Writer, errWriter io.Writer) error {
	out := make([]string, len(f.GoInputs)+3)
	out[0], out[1], out[2] = "run", f.mainFile, f.GoPipesPath
	for k, v := range f.GoInputs {
		out[k+3] = v.Path
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
		BaseDir:    baseDir,
		GoTemplate: template,
		GoInputs:   make([]struct{ Path, StructName, VarName string }, len(inputs)),
	}

	err := f.setGoInputs(inputs)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "parsing inputs", err.Error())
	}
	err = f.pipes()
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "writing pipes.go", err.Error())
	}
	err = f.mainGo()
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "creating main.go", err.Error())
	}
	return f, nil
}

func splitGoInput(input string) (fileName string, structName string) {
	//optionally the structname can be passed via commandline with ":" seperator
	if strings.Contains(input, ":") {
		s := strings.Split(input, ":")
		return s[0], s[1]
	}
	//goflat convention is to build a structname based on filename using strings Title convention
	name := filepath.Base(input)
	name = strings.Title(strings.Split(name, ".")[0])
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)
	return input, name
}

const (
	//ErrMissingOnDisk Expected error for accessing invalid file or directory
	ErrMissingOnDisk = "(file or directory is missing)"
)
