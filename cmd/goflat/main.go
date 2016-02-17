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

	b, err := goflat.NewFlat(baseDir, args.Template, args.Inputs)
	checkError(err)
	b.GoRun(os.Stdout, os.Stderr)
}

func tmpDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error %s", err.Error())
		fmt.Println()
		os.Exit(1)
	}
}
