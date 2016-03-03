package goflat_test

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/aminjam/goflat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoFlat", func() {
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

	Context("when invalid builder", func() {
		var (
			templateDir string
			builder     FlatBuilder
			buffer      bytes.Buffer
			writer      io.Writer
		)
		BeforeEach(func() {
			templateDir, _ = ioutil.TempDir(os.TempDir(), "")
			template := filepath.Join(templateDir, "test")
			err := ioutil.WriteFile(template, []byte(`Hello`), 0666)
			Expect(err).To(BeNil())

			builder, err = NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())

			writer = bufio.NewWriter(&buffer)
		})
		AfterEach(func() {
			defer os.RemoveAll(templateDir)
			buffer.Reset()
		})
		It("should catch undefined MainGo and DefaultPipes", func() {
			flat := builder.Flat()
			err := flat.GoRun(writer, writer)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMainGoUndefined))
			Expect(err.Error()).To(ContainSubstring(ErrDefaultPipesUndefined))
		})
		It("should catch undefined MainGo", func() {
			builder.EvalGoPipes("")
			flat := builder.Flat()
			err := flat.GoRun(writer, writer)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMainGoUndefined))
			Expect(err.Error()).ToNot(ContainSubstring(ErrDefaultPipesUndefined))
		})
		It("should catch undefined DefaultPipes", func() {
			builder.EvalMainGo()
			flat := builder.Flat()
			err := flat.GoRun(writer, writer)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).ToNot(ContainSubstring(ErrMainGoUndefined))
			Expect(err.Error()).To(ContainSubstring(ErrDefaultPipesUndefined))
		})
	})

	Context("when defining custom pipes", func() {
		It("should override and extend default pipes", func() {
			templateDir, _ := ioutil.TempDir(os.TempDir(), "")
			defer os.RemoveAll(templateDir)
			template := filepath.Join(templateDir, "test")
			err := ioutil.WriteFile(template, []byte(`Hello {{"oink oink oink" | replace "k" "ky"}}, tell us a {{"SECRET" | sanitize}}.`), 0666)
			Expect(err).To(BeNil())
			customPipes := filepath.Join(examples, "pipes", "pipes.go")

			builder, err := NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())

			err = builder.EvalGoPipes(customPipes)
			Expect(err).To(BeNil())
			err = builder.EvalMainGo()
			Expect(err).To(BeNil())

			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			flat := builder.Flat()
			err = flat.GoRun(writer, writer)
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(ContainSubstring("Hello oinky oink oink, tell us a TERCES."))
		})
	})

	Context("when current working direction is changed", func() {
		var wd string
		BeforeEach(func() {
			wd, _ = os.Getwd()
			err := os.Chdir(examples)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			err := os.Chdir(wd)
			Expect(err).To(BeNil())
		})
		It("should run successfully", func() {

			templateDir, _ := ioutil.TempDir(os.TempDir(), "")
			defer os.RemoveAll(templateDir)
			template := filepath.Join(templateDir, "test")
			err := ioutil.WriteFile(template, []byte(`Hello World.`), 0666)
			Expect(err).To(BeNil())

			builder, err := NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())

			err = builder.EvalGoPipes("")
			Expect(err).To(BeNil())
			err = builder.EvalMainGo()
			Expect(err).To(BeNil())

			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)
			flat := builder.Flat()
			err = flat.GoRun(writer, writer)
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(ContainSubstring("Hello World"))
		})
	})

	Context("when running the examples templates", func() {
		var (
			result      []byte
			result_file string
			buffer      bytes.Buffer

			template string
		)
		AfterEach(func() {
			buffer.Reset()
		})
		JustBeforeEach(func() {
			builder, err := NewFlatBuilder(tmpDir, template)
			Expect(err).To(BeNil())

			inputFiles := []string{
				filepath.Join(examples, "inputs", "private.go"),
				filepath.Join(examples, "inputs", "repos.go"),
			}
			err = builder.EvalGoInputs(inputFiles)
			Expect(err).To(BeNil())
			err = builder.EvalGoPipes("")
			Expect(err).To(BeNil())
			err = builder.EvalMainGo()
			Expect(err).To(BeNil())

			writer := bufio.NewWriter(&buffer)
			flat := builder.Flat()
			err = flat.GoRun(writer, writer)
			Expect(err).To(BeNil())
			result, err = ioutil.ReadFile(result_file)
			Expect(err).To(BeNil())
		})
		Describe("parsing YAML template", func() {
			BeforeEach(func() {
				template = filepath.Join(examples, "template.yml")
				result_file = filepath.Join(examples, "output.yml")
			})
			It("should show the parsed output", func() {
				Expect(result).ToNot(BeNil())
				Expect(buffer.String()).To(Equal(string(result)))
			})
		})
		Describe("parsing JSON template", func() {
			BeforeEach(func() {
				template = filepath.Join(examples, "template.json")
				result_file = filepath.Join(examples, "output.json")
			})
			It("should show the parsed output", func() {
				Expect(result).ToNot(BeNil())
				Expect(buffer.String()).To(Equal(string(result)))
			})
		})
		Describe("parsing XML template", func() {
			BeforeEach(func() {
				template = filepath.Join(examples, "template.xml")
				result_file = filepath.Join(examples, "output.xml")
			})
			It("should show the parsed output", func() {
				Expect(result).ToNot(BeNil())
				Expect(buffer.String()).To(Equal(string(result)))
			})
		})
	})
})
