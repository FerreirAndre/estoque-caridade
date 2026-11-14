package model

import "time"

type Movimentacao struct {
	ID         int64                `json:"id"`
	ItemID     int64                `json:"item_id"`
	Natureza   NaturezaMovimentacao `json:"natureza"`
	Quantidade int64                `json:"quantidade"`
	Data       time.Time            `json:"data"`
}

type NaturezaMovimentacao string

const (
	NaturezaEntrada NaturezaMovimentacao = "ENTRADA"
	NaturezaSaida   NaturezaMovimentacao = "SAIDA"
)
