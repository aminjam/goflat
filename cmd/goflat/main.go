package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aminjam/goflat"
	"github.com/jessevdk/go-flags"
)

type args struct {
	Template string   `short:"t" long:"template" description:"Template Path e.g. /PATH/TO/file.{yml,json}"`
	Inputs   []string `short:"i" long:"inputs" description:"Path to input files e.g. PATH/TO/privte.go [optional ':' struct name]"`
	Pipes    string   `short:"p" long:"pipes" description:"User defined pipes e.g. /PATH/TO/pipes.go"`
	Version  bool     `short:"v" long:"version" description:"Show version"`
}

func main() {
	args := parseArgs()
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

func parseArgs() *args {
	if len(os.Args) <= 1 {
		fmt.Println("Run --help for more help")
		os.Exit(1)
	}
	var args args
	parser := flags.NewParser(&args, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}

	if args.Version {
		fmt.Println(goflat.Version + goflat.VersionPrerelease)
		os.Exit(0)
	}
	return &args
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
