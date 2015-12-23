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
