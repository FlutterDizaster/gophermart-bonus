package models

import "time"

//go:generate easyjson -all balance.go
type Balance struct {
	Current     float64   `json:"current"`
	Withdrawn   float64   `json:"withdrawn"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

type Accrue struct {
	Username string
	Amount   float64
	OrderID  uint64
}
