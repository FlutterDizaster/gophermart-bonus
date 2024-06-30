package keygen

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

func GenerateRandomKey(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("error generating key")
	}

	return hex.EncodeToString(b)
}
