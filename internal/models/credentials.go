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

//go:generate easyjson -all credentials.go
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
