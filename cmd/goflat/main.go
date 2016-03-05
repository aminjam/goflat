package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aminjam/goflat"
	"github.com/jessevdk/go-flags"
)

type command struct {
	Template string   `short:"t" long:"template" description:"Template Path e.g. /PATH/TO/file.{yml,json}"`
	Inputs   []string `short:"i" long:"inputs" description:"Path to input files e.g. PATH/TO/privte.go [optional ':' struct name]"`
	Pipes    string   `short:"p" long:"pipes" description:"User defined pipes e.g. /PATH/TO/pipes.go"`
	Version  bool     `short:"v" long:"version" description:"Show version"`
}

func main() {
	var args command
	parser := flags.NewParser(&args, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	checkError(err)

	if args.Version {
		fmt.Println(goflat.Version + goflat.VersionPrerelease)
		os.Exit(0)
	}

	baseDir, err := tmpDir()
	if err != nil {
		checkError(fmt.Errorf("%s:%s", "cannot create temp directory", err.Error()))
	}
	defer os.RemoveAll(baseDir)

	builder, err := goflat.NewFlatBuilder(baseDir, args.Template)
	err = builder.EvalGoInputs(args.Inputs)
	checkError(err)
	err = builder.EvalGoPipes(args.Pipes)
	checkError(err)
	err = builder.EvalMainGo()
	checkError(err)

	flat := builder.Flat()
	err = flat.GoRun(os.Stdout, os.Stderr)
	checkError(err)

}

func tmpDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
