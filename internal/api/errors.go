// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package api

import "errors"

var (
	errWrongRequest       = errors.New("error wrong request")
	errUserIDNotAvaliable = errors.New("error username not avaliable")
)
