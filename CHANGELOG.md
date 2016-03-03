## 0.3.1 (03.02.2016)

### Bugfix
- This is a bug-fix that was introduced as of `v0.2.0`. Compiled binary was trying to read off of the fs for the templates and default pipes. By using `go generate` we can still test the behavior of the extensions while embedding the text into the compiled binary.

## 0.3.0 (03.01.2016)

### Features
- Enable user defined pipes e.g. `--pipes FILE.go`
- A user should be able to override or extend default helper functions

## 0.2.0 (02.25.2016)

### Features
- Add support for pipes and helper functions e.g. `{{.List | toUpper}}`
- Support toUpper, toLower, split, join, replace, map

### Breaking Changes
- [Changed implementaion](https://github.com/aminjam/goflat/commit/89a00c8abb54e341f935ff4547da382ff4efa51f)
 for `{{join .List ","}}` to `{{.List | join ","}}`
