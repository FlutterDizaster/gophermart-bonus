package models

import "time"

type OrderStatus string

const (
	StatusOrderNew        = "NEW"
	StatusOrderProcessing = "PROCESSING"
	StatusOrderInvalid    = "INVALID"
	StatusOrderProcessed  = "PROCESSED"
)

//easyjson:json
type Orders []Order

//go:generate easyjson -all order.go
type Order struct {
	ID         uint64      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
