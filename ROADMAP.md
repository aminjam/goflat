## 0.4.0
- Adding `--output` option for writing to a file.
- Adding `go get` support for missing imports. For example if `gopkg.in/yaml.v2` is used within an input and not in `GOPATH`, goflat should `go get` the missing dependencies.
