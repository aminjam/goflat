package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/aminjam/goflat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoFlat", func() {

	Context("With having invalid params", func() {
		It("should catch invalid baseDir", func() {
			b, err := NewFlat("INVALID", "INVALID", nil)
			Expect(b).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
		})
		It("should catch invalid input files", func() {
			tmpDir, _ := ioutil.TempDir(os.TempDir(), "")
			defer os.RemoveAll(tmpDir)
			wd, _ := os.Getwd()
			template := filepath.Join(wd, "fixtures", "pipeline.yml")

			invalid_input_files := []string{"/WRONG/FILE", "WRONG/ANOTHER/FILES"}

			b, err := NewFlat(tmpDir, template, invalid_input_files)
			Expect(b).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
			Expect(err.Error()).To(ContainSubstring("/WRONG/FILE"))
		})

		It("should catch invalid template", func() {
			tmpDir, _ := ioutil.TempDir(os.TempDir(), "")
			defer os.RemoveAll(tmpDir)

			invalid_template := "/WRONG/FILE.yml"

			b, err := NewFlat(tmpDir, invalid_template, []string{})
			Expect(b).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring(ErrMissingOnDisk))
			Expect(err.Error()).To(ContainSubstring(invalid_template))
		})
	})
	Context("", func() {
		wd, _ := os.Getwd()
		template := filepath.Join(wd, "fixtures", "pipeline.yml")
		inputFiles := []string{
			filepath.Join(wd, "fixtures", "private.go") + ":Private",
			filepath.Join(wd, "fixtures", "repos.go"),
		}

		Describe("When input is not from a .go extension file", func() {
			It("should successfully copy the file with an attached extension", func() {
				tmpDir, _ := ioutil.TempDir(os.TempDir(), "")
				defer os.RemoveAll(tmpDir)
				orgFile := filepath.Join(wd, "fixtures", "a-private-note")
				b, err := NewFlat(tmpDir, template, []string{
					orgFile + ":Private",
				})
				Expect(err).To(BeNil())

				var newFile string
				for _, v := range b.Inputs {
					if filepath.Base(v.Path) == "a-private-note.go" {
						newFile = v.Path
						break
					}
				}
				Expect(newFile).ToNot(BeNil())
				newFileInfo, err := os.Stat(newFile)
				Expect(err).To(BeNil())
				orgFileInfo, _ := os.Stat(orgFile)
				Expect(orgFileInfo.Size()).To(Equal(newFileInfo.Size()))
			})
		})

		Describe("When valid params", func() {
			var (
				tmpDir string
				flat   *Flat
			)
			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir(os.TempDir(), "")
				b, err := NewFlat(tmpDir, template, inputFiles)
				Expect(err).To(BeNil())
				flat = b
			})
			AfterEach(func() {
				defer os.RemoveAll(tmpDir)
			})

			It("should successfully copy a valid file", func() {
				orgFile := filepath.Join(wd, "fixtures", "repos.go")
				var newFile string
				for _, v := range flat.Inputs {
					if filepath.Base(v.Path) == "repos.go" {
						newFile = v.Path
						break
					}
				}
				Expect(newFile).ToNot(BeNil())
				newFileInfo, err := os.Stat(newFile)
				Expect(err).To(BeNil())
				orgFileInfo, _ := os.Stat(orgFile)
				Expect(orgFileInfo.Size()).To(Equal(newFileInfo.Size()))
			})
			It("should have created main.go ", func() {
				mainGoFile := filepath.Join(tmpDir, "main.go")
				fileInfo, err := os.Stat(mainGoFile)
				Expect(err).To(BeNil())
				Expect(fileInfo).ToNot(BeNil())

				data, err := ioutil.ReadFile(mainGoFile)
				Expect(err).To(BeNil())
				Expect(data).To(ContainSubstring(fmt.Sprintf("data, err := ioutil.ReadFile(\"%s\")", flat.Template)))
				Expect(data).To(ContainSubstring(fmt.Sprintf(
					"result.%s = New%s()", flat.Inputs[0].StructName, flat.Inputs[0].StructName)))
			})
			It("should output the parsed template ", func() {
				var b bytes.Buffer
				writer := bufio.NewWriter(&b)
				err := flat.GoRun(writer, writer)
				Expect(err).To(BeNil())
				// copy the output in a separate goroutine so printing can't block indefinitely
				result_file := filepath.Join(wd, "fixtures", "result.yml")
				result, err := ioutil.ReadFile(result_file)
				Expect(err).To(BeNil())
				Expect(result).ToNot(BeNil())
				Expect(b.String()).To(Equal(string(result)))
			})
		})

	})
})
