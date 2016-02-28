package goflat_test

import (
	"bufio"
	"bytes"
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

	Context("when running the examples templates", func() {
		var (
			result      []byte
			template    string
			result_file string
			buffer      bytes.Buffer
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
			err = builder.EvalGoPipes()
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
