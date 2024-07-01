// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package api

import (
	"errors"
	"net/http"
	"time"

	sharederrors "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

// Handler для регистрации пользователей.
func (api *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных для регистрации пользователя
	cred, err := getCredentialsFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Регистрация пользователя
	token, err := api.userMgr.Register(r.Context(), cred)
	if err != nil {
		if errors.Is(err, sharederrors.ErrUserAlreadyExist) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Создание cookie
	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(api.cookieTTL),
	}

	// Установка cookie
	http.SetCookie(w, cookie)

	w.Header().Add("Authorization", token)

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}
