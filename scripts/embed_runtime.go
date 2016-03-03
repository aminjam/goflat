package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fs, _ := ioutil.ReadDir("runtime")
	runtime_files := []string{"main.gotempl", "pipes.go"}

	out, _ := os.Create("runtime.go")
	out.Write([]byte("package goflat\n\nconst (\n"))
	for _, f := range fs {
		for _, runtime_file := range runtime_files {
			if runtime_file == f.Name() {
				varName := strings.Replace(strings.Title(f.Name()), ".", "", -1)
				fPath := filepath.Join("runtime", f.Name())
				out.Write([]byte(varName + " = `"))
				f, _ := os.Open(fPath)
				io.Copy(out, f)
				out.Write([]byte("`\n"))
			}
		}
	}
	out.Write([]byte(")\n"))
}
