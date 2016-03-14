package goflat

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:generate go run scripts/embed_runtime.go

//Flat struct
type Flat struct {
	MainGo       string
	GoTemplate   string
	GoInputs     []goInput
	DefaultPipes string
	CustomPipes  string

	goPath string
}

//GoRun runs go on the dynamically created main.go with a given stdout and stderr pipe
func (f *Flat) GoRun(outWriter io.Writer, errWriter io.Writer) error {
	err := f.validate()
	if err != nil {
		return err
	}
	err = f.goGetImports()
	if err != nil {
		return err
	}
	out := []string{"run", f.MainGo, f.DefaultPipes}

	if f.CustomPipes != "" {
		out = append(out, f.CustomPipes)
	}
	for _, v := range f.GoInputs {
		out = append(out, v.Path)
	}
	cmd := exec.Command("go", out...)
	cmd.Env = f.env()
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	return cmd.Run()
}

func (f *Flat) goGetImports() error {
	out := []string{"get", "./..."}
	cmd := exec.Command("go", out...)
	cmd.Dir = f.goPath
	cmd.Env = f.env()

	writer := bytes.NewBufferString("")
	cmd.Stdout = writer
	cmd.Stderr = writer
	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s", err.Error(), writer.String()))
	}
	return nil
}

func (f *Flat) env() []string {
	env := environ(os.Environ())
	old_gopath := os.Getenv("GOPATH")
	if old_gopath != "" {
		old_gopath = ":" + old_gopath
	}
	env.Unset("GOPATH")
	gopath := fmt.Sprintf("GOPATH=%s%s", f.goPath, old_gopath)
	return append(env, gopath)
}

func (f *Flat) validate() error {
	msgs := []string{}
	if f.MainGo == "" {
		msgs = append(msgs, ErrMainGoUndefined)
	}
	if f.DefaultPipes == "" {
		msgs = append(msgs, ErrDefaultPipesUndefined)
	}

	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, ","))
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type goInput struct{ Path, StructName, VarName string }

func newGoInput(input string) goInput {
	//optionally the structname can be passed via commandline with ":" seperator
	if strings.Contains(input, ":") {
		s := strings.Split(input, ":")
		return goInput{Path: s[0], StructName: s[1]}
	}
	//goflat convention is to build a structname based on filename using strings Title convention
	name := filepath.Base(input)
	name = strings.Title(strings.Split(name, ".")[0])
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)
	return goInput{
		Path:       input,
		StructName: name,
	}
}

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

// Unset a single environment variable.
func (e *environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

const (
	ErrMainGoUndefined       = "(main func is missing)"
	ErrDefaultPipesUndefined = "(default pipes file is missing)"
)
