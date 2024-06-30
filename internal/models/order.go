// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package models

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
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt string      `json:"uploaded_at"`
}
