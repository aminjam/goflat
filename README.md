# goflat [![Build Status](https://travis-ci.org/aminjam/goflat.png?branch=master)](https://travis-ci.org/aminjam/goflat)
A Go template flattener `goflat` is for creating complex configuration files (JSON, YAML, XML, etc.).

## Motivation
Building long configuration files is not fun nor testable. Replacing passwords and secrets in a configuration file is usually done with regex and sometimes it's unpredictable! Why not use go templates, along with individual `.go` input files, that know how to unmarshall and parse their own data structure. This way we can build a complex configuration with inputs coming from different `structs` that we can test their behavior independently. That is what `goflat` does. A small and simple go template flattener that uses go runtime to dynamically create a template for parsing go structs.

## Getting Started

### Run as executable
```
go get github.com/aminjam/goflat/cmd/goflat
$GOPATH/bin/goflat --help
```
### Run from source
Built with Go 1.5.3 and `GO15VENDOREXPERIMENT` flag.
```
git clone https://github.com/aminjam/goflat.git && cd goflat
make init
make build
./pkg/*/goflat --help
```
### Usage
```
goflat -t FILE.{yml,json,xml} -i private.go -i teams.go -i repos.go ...
```
```
goflat -t FILE.{yml,json,xml} -i <(lpass show 'private.go' --notes):Private
```
## Example

Here is a sample YAML configuration used for creating [concourse](https://concourse.ci) pipeline.
```
{{ $global := . }}
resources:
- name: ci
  type: git
  source:
    uri: https://github.com/cloudfoundry/myproject-ci.git
{{range .Repos}}
- name: {{.Name}}
  type: git
  source:
    uri: {{.Repo}}
    branch: {{.Branch}}
{{end}}

jobs:
{{range .Repos}}
- name: {{.Name}}
  serial: true
  plan:
  - aggregate:
    - get: ci
    - get: project
      resource: {{.Name}}
      trigger: true
  - task: do-something
    config:
      platform: linux
      image: "docker:///alpine"
      run:
        path: sh
        args: ["-c","echo Hi"]
      params:
        PASSWORD: {{$global.Private.Password}}
        SECRET: {{$global.Private.Secret}}
{{end}}
- name-{{.Private.Secret}}: {{.Private.Password}}
- comma-seperated-repo-names: {{.Repos.Names | join ","}}
```
We have `Repos` and `Private` struct that contain some runtime data that needs to be parsed into the template. Here is a look at the checked-in `private.go`

```
package main

type Private struct {
	Password string
	Secret   string
}

func NewPrivate() Private {
	return Private{
		Password: "team3",
		Secret:   "cloud-foundry",
	}
}
```
Each of the input files are required to have 2 things:
* A struct named after the filename (e.g. filename `hello-world.go` should have `HelloWorld` struct). If the struct name differs from the filename convention, you can optionally provide the name of the struct (e.g. `-i <(lpass show 'file.go' --notes):Private`)
* A `New{{.StructName}}` function that returns `{{.StructName}}` (e.g. `func NewPrivate() Private{}`)

Similarly, we can also define `repos.go` as an array of objects to use within `{{range .Repos}}`.
```
package main

type Repos []struct {
	Name   string
	Repo   string
	Branch string
}

func (r Repos) Names() []string {
	names := make([]string, len(r))
	for k, v := range r {
		names[k] = v.Name
	}
	return names
}

func NewRepos() Repos {
	return Repos{
		{
			Name:   "repo1",
			Repo:   "https://github.com/jane/repo1",
			Branch: "master",
		},
		{
			Name:   "repo2",
			Repo:   "https://github.com/john/repo2",
			Branch: "develop",
		},
	}
}
```
Now we can run the sample configuration in the `.examples`.

```
goflat -t .examples/template.yml -i .examples/inputs/repos.go -i .examples/inputs/private.go
```

### Pipes "|"
Pipes can be nested and here is a set of supported helper functions:

- **join**: `{{.List | join "," }}`
- **map**: `{{.ListOfObjects | map "Name,Age" "," }}` (comma seperated property names)
- **replace**: `{{.StringValue | replace "," " " }}`
- **split**: `{{.StringValue | split "," }}`
- **toLower**: `{{.Field | toLower }}`
- **toUpper**: `{{.Field | toUpper }}`

You can optionally define a custom list of helper functions that overrides or extends the behavior of the default pipes. See [an exmaple](.examples/pipes/pipes.go) file that can optionally be passed via `--pipes` flag. Note that the function signature has to be the following:

```
func CustomPipes() tempate.FuncMap {
...
}
```
