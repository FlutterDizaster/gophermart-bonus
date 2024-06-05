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
	ErrUserAlreadyExist     = errors.New("error user already exist")
	ErrWrongLoginOrPassword = errors.New("error wrong login or password")
)
