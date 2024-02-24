package main

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type Transaction struct {
	Valor       int        `json:"valor"`
	Tipo        string     `json:"tipo"`
	Descricao   string     `json:"descricao"`
	RealizadaEm *time.Time `json:"realizada_em"`
}

type Client struct {
	ID     int `json:"id"`
	Total  int `json:"total"`
	Limite int `json:"limite"`
}
