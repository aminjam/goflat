package runtime_test

import (
	"encoding/base64"
	"errors"
	"strings"
	"text/template"

	. "github.com/aminjam/goflat/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Pipes", func() {
	var (
		pipes  *Pipes
		tmpl   *template.Template
		buffer *gbytes.Buffer
	)
	BeforeEach(func() {
		pipes = NewPipes()
		tmpl = template.New("tester").Funcs(pipes.Map)
		buffer = gbytes.NewBuffer()
	})

	Describe("when extending default pipes", func() {
		JustBeforeEach(func() {
			fm := template.FuncMap{
				"join": func(a []string) (string, error) {
					return strings.Join(a, ":-)"), nil
				},
				"base64": func(method string, a string) (string, error) {
					if method == "encode" {
						return base64.StdEncoding.EncodeToString([]byte(a)), nil
					} else if method == "decode" {
						out, err := base64.StdEncoding.DecodeString(a)
						return string(out), err
					} else {
						return "", errors.New("method not found")
					}
				},
			}
			pipes.Extend(fm)
		})
		It("should add base64 method", func() {
			const text = `{{ . | base64 "encode" }}`
			tmpl := template.New("tester").Funcs(pipes.Map)
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			buffer := gbytes.NewBuffer()
			err = tmpl.Execute(buffer, "abcd")
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`YWJjZA==`))
		})
		It("should override the joins method", func() {
			const text = `{{ . | join }}`
			tmpl := template.New("tester").Funcs(pipes.Map)
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			buffer := gbytes.NewBuffer()
			err = tmpl.Execute(buffer, []string{"a", "b"})
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`a\:\-\)b`))
		})
	})

	Describe("when using default pipes", func() {
		It("should validate joins method", func() {
			const text = `{{ . | join "," }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, []string{"a", "b"})
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`a,b`))
		})
		Context("when validating a map method", func() {
			type Tester struct{ Name, Job string }
			t := []Tester{
				{"John", "Jabber"},
				{"Cherry", "Chatter"},
			}

			It("should accepet a single Field", func() {
				const text = `{{ . | map "Job" "|" }}`
				tmpl, err := tmpl.Parse(text)
				Expect(err).To(BeNil())
				err = tmpl.Execute(buffer, t)
				Expect(err).To(BeNil())
				Eventually(buffer).Should(gbytes.Say("Jabber Chatter"))
			})
			It("should accepet a comma seperated fields", func() {
				const text = `{{ . | map "Job,Name" "." }}`
				tmpl, err := tmpl.Parse(text)
				Expect(err).To(BeNil())
				err = tmpl.Execute(buffer, t)
				Expect(err).To(BeNil())
				Eventually(buffer).Should(gbytes.Say("Jabber.John Chatter.Cherry"))
			})
		})
		It("should validate replace method", func() {
			const text = `{{ . | replace "A" "D" }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "AbCDAa")
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`DbCDDa`))
		})
		It("should validate split method", func() {
			const text = `{{ range (. | split " ") }}{{.}}-item {{end}}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "AB CD")
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`AB-item CD-item `))
		})
		It("should validate toUpper method", func() {
			const text = `{{ . | toUpper }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "abA")
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`ABA`))
		})
		It("should validate toLower method", func() {
			const text = `{{ . | toLower }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "BaBA")
			Expect(err).To(BeNil())
			Eventually(buffer).Should(gbytes.Say(`baba`))
		})
	})
})
