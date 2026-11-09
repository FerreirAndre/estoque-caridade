package model

type Item struct {
	ID         int64         `json:"id"`
	Nome       string        `json:"nome"`
	Quantidade int64         `json:"quantidade"`
	Unidade    UnidadeMedida `json:"unidade_medida"`
}
