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
