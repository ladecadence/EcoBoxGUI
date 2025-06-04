package models

type Container struct {
	Code   string `json:"codigo"`
	NameES string `json:"nombre_es"`
	NameEU string `json:"nombre_eu"`
	Active int    `json:"active"`
	Icon   string `json:"icono"`
	Price  string `json:"precio"`
}

type Containers struct {
	Response
	Containers []Container `json:"contenedores"`
}
