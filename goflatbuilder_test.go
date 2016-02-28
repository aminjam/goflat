package goflat_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/aminjam/goflat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoFlatBuilder", func() {
	var (
		tmpDir   string
		examples string
	)
	BeforeEach(func() {
		tmpDir, _ = ioutil.TempDir(os.TempDir(), "")
		wd, _ := os.Getwd()
		examples = filepath.Join(wd, ".examples")
	})
	AfterEach(func() {
		defer os.RemoveAll(tmpDir)
	})

	Context("with invalid params", func() {
		It("should catch invalid baseDir", func() {
			builder, err := NewFlatBuilder("INVALID", "INVALID")
			Expect(builder).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
		})
		It("should catch invalid input files", func() {
			template := filepath.Join(examples, "template.yml")
			invalid_input_files := []string{"/WRONG/FILE", "WRONG/ANOTHER/FILES"}

			builder, err := NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())
			err = builder.EvalGoInputs(invalid_input_files)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
			Expect(err.Error()).To(ContainSubstring("/WRONG/FILE"))
		})

		It("should catch invalid template", func() {
			invalid_template := "/WRONG/FILE.yml"

			builder, err := NewFlatBuilder(tmpDir, invalid_template)
			Expect(builder).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
			Expect(err.Error()).To(ContainSubstring(invalid_template))
		})
	})
	Context("#EvalGoInputs", func() {
		var builder FlatBuilder
		BeforeEach(func() {
			var err error
			template := filepath.Join(examples, "template.yml")
			builder, err = NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())
		})
		It("should create destination input file", func() {
			inputFiles := []string{
				filepath.Join(examples, "inputs", "private.go"),
			}
			err := builder.EvalGoInputs(inputFiles)
			Expect(err).To(BeNil())
			flat := builder.Flat()
			Expect(len(flat.GoInputs)).To(Equal(1))
			newFileInfo, err := os.Stat(flat.GoInputs[0].Path)
			Expect(err).To(BeNil())
			orgFileInfo, _ := os.Stat(inputFiles[0])
			Expect(orgFileInfo.Size()).To(Equal(newFileInfo.Size()))
		})
		It("should evaluate files with custom struct", func() {
			orgFile := filepath.Join(examples, "inputs", "a-private-note")
			err := builder.EvalGoInputs([]string{
				orgFile + ":Private",
			})
			Expect(err).To(BeNil())

			flat := builder.Flat()
			Expect(len(flat.GoInputs)).To(Equal(1))
			newFileInfo, err := os.Stat(flat.GoInputs[0].Path)
			Expect(err).To(BeNil())
			orgFileInfo, _ := os.Stat(orgFile)
			Expect(orgFileInfo.Size()).To(Equal(newFileInfo.Size()))
		})
	})
	Context("#EvalGoPipes", func() {
		var builder FlatBuilder
		BeforeEach(func() {
			var err error
			template := filepath.Join(examples, "template.yml")
			builder, err = NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())
		})
		It("should create destination input file", func() {
			pipesFile := filepath.Join(examples, "pipes", "pipes.go")
			err := builder.EvalGoPipes(pipesFile)
			Expect(err).To(BeNil())
			flat := builder.Flat()
			Expect(flat.CustomPipes).ToNot(BeEmpty())

			newFileInfo, err := os.Stat(flat.CustomPipes)
			Expect(err).To(BeNil())
			orgFileInfo, _ := os.Stat(pipesFile)
			Expect(orgFileInfo.Size()).To(Equal(newFileInfo.Size()))
		})
	})
	Context("#EvalMainGo", func() {
		It("should have created main.go", func() {
			var (
				templateDir, template, infoGo string
				err                           error
			)
			templateDir, _ = ioutil.TempDir(os.TempDir(), "")
			defer os.RemoveAll(templateDir)
			template = filepath.Join(templateDir, "test.txt")
			err = ioutil.WriteFile(template, []byte("Hello {{.Info.Name}}"), 0666)
			Expect(err).To(BeNil())
			infoGo = filepath.Join(templateDir, "info.go")
			err = ioutil.WriteFile(infoGo, []byte(`package main
			type Info struct { Name string }
			func NewInfo() Info { return Info { Name: "Jane" } }`), 0666)
			Expect(err).To(BeNil())

			builder, err := NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())
			err = builder.EvalGoInputs([]string{infoGo})
			Expect(err).To(BeNil())
			err = builder.EvalMainGo()
			Expect(err).To(BeNil())
			flat := builder.Flat()

			newFileInfo, err := os.Stat(flat.MainGo)
			Expect(err).To(BeNil())
			Expect(newFileInfo).ToNot(BeNil())

			data, err := ioutil.ReadFile(flat.MainGo)
			Expect(err).To(BeNil())
			Expect(data).To(ContainSubstring(fmt.Sprintf("data, err := ioutil.ReadFile(\"%s\")", flat.GoTemplate)))
			Expect(data).To(ContainSubstring(fmt.Sprintf(
				"result.%s = New%s()", flat.GoInputs[0].StructName, flat.GoInputs[0].StructName)))
		})
	})
})
