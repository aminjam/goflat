# go-flat
A Go template flattener `go-flat` is for creating complex configuration files (JSON, YAML, etc.) with secrets.

## Motivation
Building long YAML or JSON files is not fun! Replacing passwords and secrets in a configuration file is usually done with regex and it's unreliable! Why not use go templates, along with individual `.go` input files, that know how to unmarshall and parse their own data structure?! This way we can build a complex configuration file with inputs coming from different `structs`. That is what `go-flat` does. A small and simple go template flattener uses go runtime to dynamically create a template for parsing go structs.

## Getting Started

### Run as executable
```
go get github.com/aminjam/go-flat
$GOPATH/bin/go-flat --help
```
### Run from source
```
git clone https://github.com/aminjam/go-flat.git && cd go-flat
make update-deps
make build
./pkg/*/go-flat --help
```
### Usage
```
go-flat -t FILE.{yml,json} -i private.go -i teams.go -i repos.go ...
```
```
go-flat -t FILE.{yml,json} -i <(lpass show 'private.go' --notes):Private
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
```
We Have `Repos` and `Private` struct that contain some runtime data that needs to be parsed into the template. Here is a look at the checked-in `private.go`

```
package main

import "encoding/json"

type Private struct {
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type privateAlias struct {
	Data    string
	Private Private
}

func (rs *privateAlias) Flat() (Private, error) {
	data := []byte(rs.Data)
	err := json.Unmarshal(data, &rs)
	return rs.Private, err
}

func NewPrivate() *privateAlias {
	return &privateAlias{
		Data: `
{
	"private":{
		"password":"team3",
		"secret":"cloud-foundry"
	}
}
`,
		Private: Private{},
	}
}
```
Each of the input files are required to have 4 things:
* A struct named after the filename (e.g. filename `hello-world.go` should have `HelloWorld` struct). If the struct name differs from the filename convention, you can optionally provide the name of the struct (e.g. `-i <(lpass show 'file.go' --notes):Private`)
* A `New{{.StructName}}` function
* An Unmarshaller `Flat()` function that serializes the object
* A package name should always be `main`

Similarly, we can also define `repos.go` as an array of objects to use within `{{range .Repos}}`
```
package main

import "gopkg.in/yaml.v2"

type Repos []struct {
	Name   string `yaml:"name"`
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
}

type ReposStructure struct {
	Data  string
	Repos Repos
}

func (rs *ReposStructure) Flat() (Repos, error) {
	data := []byte(rs.Data)
	err := yaml.Unmarshal(data, &rs)
	return rs.Repos, err
}

func NewRepos() *ReposStructure {
	return &ReposStructure{
		Data: `
repos:
- name: repo1
  repo: https://github.com/jane/repo1
  branch: master
- name: repo2
  repo: https://github.com/john/repo2
  branch: develop
`,
		Repos: make(Repos, 2),
	}
}
```
Note that each of the inputs can have their own serialization method, as long as the serialization library is in the main `$GOPATH`. Now we can run the example pipeline in the fixtures.
```
goflat -t fixtures/pipeline.yml -i fixtures/repos.go -i fixtures/private.go
```
