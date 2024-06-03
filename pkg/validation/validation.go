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
