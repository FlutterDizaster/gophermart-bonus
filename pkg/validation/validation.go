// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package validation

import (
	"crypto/sha256"
)

func CalculateHashSHA256WithKey(content, key []byte) []byte {
	content = append(content, key...)
	return CalculateHashSHA256(content)
}

func CalculateHashSHA256(content []byte) []byte {
	hash := sha256.Sum256(content)
	return hash[:]
}
