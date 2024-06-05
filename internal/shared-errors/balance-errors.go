// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package sharederrors

import "errors"

var (
	ErrNotEnoughFunds     = errors.New("error not enougs funds")
	ErrWithdrawNotAllowed = errors.New("error withdraw not allowed")
	ErrWrongOrderID       = errors.New("error wrong order id")
)
