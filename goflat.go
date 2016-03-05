package goflat

import (
	"errors"
	"io"
	"math/rand"
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
}

//GoRun runs go on the dynamically created main.go with a given stdout and stderr pipe
func (f *Flat) GoRun(outWriter io.Writer, errWriter io.Writer) error {
	err := f.validate()
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
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	return cmd.Run()
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

const (
	ErrMainGoUndefined       = "(main func is missing)"
	ErrDefaultPipesUndefined = "(default pipes file is missing)"
)
