package goflat_test

import (
	"text/template"

	. "github.com/aminjam/goflat"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Pipes", func() {
	var (
		pipes *Pipes
	)
	BeforeEach(func() {
		pipes = NewPipes()
	})

	Describe("when testing default pipes", func() {
		var (
			tmpl   *template.Template
			buffer *gbytes.Buffer
		)
		BeforeEach(func() {
			tmpl = template.New("tester").Funcs(pipes.Map)
			buffer = gbytes.NewBuffer()
		})
		It("should validate joins method", func() {
			const text = `{{ . | join "," }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, []string{"a", "b"})
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
			Eventually(buffer).Should(gbytes.Say(`DbCDDa`))
		})
		It("should validate split method", func() {
			const text = `{{ range (. | split " ") }}{{.}}-item {{end}}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "AB CD")
			Eventually(buffer).Should(gbytes.Say(`AB-item CD-item `))
		})
		It("should validate toUpper method", func() {
			const text = `{{ . | toUpper }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "abA")
			Eventually(buffer).Should(gbytes.Say(`ABA`))
		})
		It("should validate toLower method", func() {
			const text = `{{ . | toLower }}`
			tmpl, err := tmpl.Parse(text)
			Expect(err).To(BeNil())
			err = tmpl.Execute(buffer, "BaBA")
			Eventually(buffer).Should(gbytes.Say(`baba`))
		})
	})
})
