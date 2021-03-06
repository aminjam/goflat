package goflat

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//Builder pattern seems to be the most appropriate structure for building a `Flat` struct
type FlatBuilder interface {
	EvalGoInputs(files []string) error
	EvalGoPipes(file string) error
	EvalMainGo() error
	Flat() *Flat
}

type flatBuilder struct {
	flat    *Flat
	baseDir string
}

func (builder *flatBuilder) cp(file string) (string, error) {
	outFile := filepath.Join(builder.baseDir, nameGenerator())
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
		builder.flat.GoInputs[k] = gi
	}
	return nil
}

func (builder *flatBuilder) EvalGoPipes(file string) error {
	defaultPipes, err := builder.defaultPipes()
	if err != nil {
		return err
	}
	builder.flat.DefaultPipes = defaultPipes

	if file != "" {
		customPipes, err := builder.cp(file)
		if err != nil {
			return fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
		}
		builder.flat.CustomPipes = customPipes
	}
	return nil
}

func (builder *flatBuilder) EvalMainGo() error {
	outFile := filepath.Join(builder.baseDir, nameGenerator())
	main, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer main.Close()

	var tmpl = template.Must(template.New("main").Parse(MainGotempl))
	if err := tmpl.Execute(main, builder.flat); err != nil {
		return err
	}
	builder.flat.MainGo = outFile
	return nil
}

func (builder *flatBuilder) Flat() *Flat {
	return builder.flat
}

//NewFlatBuilder initializes a new instance of `FlatBuilder` interface
func NewFlatBuilder(baseDir, template string) (FlatBuilder, error) {
	if _, err := os.Stat(baseDir); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}
	if _, err := os.Stat(template); err != nil {
		return nil, fmt.Errorf("%s:%s", ErrMissingOnDisk, err.Error())
	}

	src_dir := filepath.Join(baseDir, "src")
	os.MkdirAll(src_dir, 0777)

	goflatDir, _ := ioutil.TempDir(src_dir, "goflat")
	goPath, _ := filepath.Abs(baseDir)
	builder := &flatBuilder{
		baseDir: goflatDir,
		flat: &Flat{
			GoTemplate: template,
			goPath:     goPath,
		},
	}

	return builder, nil
}

func (builder *flatBuilder) defaultPipes() (string, error) {
	content := strings.Replace(PipesGo, "package runtime", "package main", -1)
	outFile := filepath.Join(builder.baseDir, nameGenerator())
	err := ioutil.WriteFile(outFile, []byte(content), 0666)
	if err != nil {
		return "", err
	}
	return outFile, nil
}

func nameGenerator() string {
	var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	size := 10
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf) + ".go"
}

const (
	//ErrMissingOnDisk Expected error for accessing invalid file or directory
	ErrMissingOnDisk = "(file or directory is missing)"
)
