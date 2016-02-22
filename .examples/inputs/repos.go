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
