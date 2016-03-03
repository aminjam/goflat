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
