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
