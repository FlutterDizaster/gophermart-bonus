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
	ErrOrderAlreadyLoaded      = errors.New("error order already loaded")
	ErrOrderLoadedByAnotherUsr = errors.New("error order loaded by another user")
	ErrNoOrdersFound           = errors.New("error no orders found")
)
