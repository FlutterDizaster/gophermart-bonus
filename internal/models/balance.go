// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License.
//
// This file uses a third-party package easyjson licensed under MIT License.
//
// See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package models

//go:generate easyjson -all balance.go
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Accrue struct {
	Username string
	Amount   float64
	OrderID  uint64
}
