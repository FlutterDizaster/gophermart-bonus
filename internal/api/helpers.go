// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package api

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	ctxkeys "github.com/FlutterDizaster/gophermart-bonus/internal/context-keys"
	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
)

func getCredentialsFromReq(r *http.Request) (models.Credentials, error) {
	var cred models.Credentials
	var buf bytes.Buffer // TODO: Получение буфера из пула

	// Проверка Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return cred, errWrongRequest
	}

	// Чтение тела запроса
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return cred, errWrongRequest
	}

	// Unmarshal тела запроса
	if err := cred.UnmarshalJSON(buf.Bytes()); err != nil {
		return cred, errWrongRequest
	}

	return cred, nil
}

func getWithdrawFromReq(r *http.Request) (models.Withdraw, error) {
	var withdraw models.Withdraw
	var buf bytes.Buffer

	// Проверка Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return withdraw, errWrongRequest
	}

	// Чтение тела запроса
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return withdraw, errWrongRequest
	}

	// Unmarshal тела запроса
	if err := withdraw.UnmarshalJSON(buf.Bytes()); err != nil {
		return withdraw, errWrongRequest
	}

	return withdraw, nil
}

func getUserIDFromReq(r *http.Request) (uint64, error) {
	// Получение userID
	reqCtx := r.Context()
	userIDRaw := reqCtx.Value(ctxkeys.UserID)
	if userIDRaw == nil {
		return 0, errUserIDNotAvaliable
	}
	userID, ok := userIDRaw.(uint64)
	if !ok {
		return 0, errUserIDNotAvaliable
	}
	return userID, nil
}

// Имплементация алгоритма Луна.
func checkLuhn(number string) bool {
	number = strings.TrimSpace(number)
	seq := strings.Split(number, "")
	var sum int
	seqLen := len(seq)
	parity := seqLen % 2

	for i := 0; i < seqLen; i++ {
		digit, err := strconv.Atoi(seq[i])
		if err != nil {
			return false
		}

		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
