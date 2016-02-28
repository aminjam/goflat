package goflat

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FlatBuilder interface {
	EvalGoInputs(files []string) error
	EvalGoPipes() error
	EvalMainGo() error
	Flat() *Flat
}

type flatBuilder struct {
	flat    *Flat
	baseDir string
}

func (builder *flatBuilder) cp(file string) (string, error) {
	base := filepath.Base(file)
	if !strings.HasSuffix(base, ".go") {
		base += ".go"
	}
	outFile := filepath.Join(builder.baseDir, base)
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

func (builder *flatBuilder) EvalGoInputs(files []string) error {
	builder.flat.GoInputs = make([]goInput, len(files))
	for k, v := range files {
		gi := newGoInput(v)
		file, err := builder.cp(gi.Path)
		if err != nil {
			return fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
		}
		gi.Path = file
		gi.VarName = strings.ToLower(gi.StructName)
		builder.flat.GoInputs[k] = gi
	}
	return nil
}

func (builder *flatBuilder) EvalGoPipes() error {
	data, err := ioutil.ReadFile("pipes.go")
	if err != nil {
		return err
	}
	content := strings.Replace(string(data), "package goflat", "package main", -1)
	outFile := filepath.Join(builder.baseDir, "pipes.go")
	err = ioutil.WriteFile(outFile, []byte(content), 0666)
	if err != nil {
		return err
	}
	builder.flat.GoPipes = outFile
	return nil
}

func (builder *flatBuilder) EvalMainGo() error {
	outFile := filepath.Join(builder.baseDir, "main.go")
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
	if err := tmpl.Execute(main, builder.flat); err != nil {
		return err
	}
	builder.flat.MainGo = outFile
	return nil
}

func (builder *flatBuilder) Flat() *Flat {
	return builder.flat
}

//NewFlatBuilder initializes a new instance of Flat builder interface
func NewFlatBuilder(baseDir, template string) (FlatBuilder, error) {
	if _, err := os.Stat(baseDir); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}
	if _, err := os.Stat(template); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}
	builder := &flatBuilder{
		baseDir: baseDir,
		flat: &Flat{
			GoTemplate: template,
		},
	}

	return builder, nil
}

const (
	//ErrMissingOnDisk Expected error for accessing invalid file or directory
	ErrMissingOnDisk = "(file or directory is missing)"
)
