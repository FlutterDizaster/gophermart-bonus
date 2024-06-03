package models

//easyjson:json
type Withdraws []Withdraw

//go:generate easyjson -all withdraw.go
type Withdraw struct {
	OrderID uint64  `json:"order"`
	Sum     float64 `json:"sum"`
}
