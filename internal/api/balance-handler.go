// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package api

import (
	"log/slog"
	"net/http"
)

func (api *API) balanceHandler(w http.ResponseWriter, r *http.Request) {
	// Получение имени пользователя
	userID, err := getUserIDFromReq(r)
	if err != nil {
		slog.Error(
			"balance manager error",
			slog.String("method", "getUserIDFromReq"),
			slog.Any("error", err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получение баланса пользователя
	balance, err := api.BalanceMgr.Get(r.Context(), userID)
	if err != nil {
		slog.Error(
			"balance manager error",
			slog.String("method", "api.BalanceMgr.Get"),
			slog.Any("error", err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal ответа
	body, err := balance.MarshalJSON()
	if err != nil {
		slog.Error(
			"balance manager error",
			slog.String("method", "MarshalJSON"),
			slog.Any("error", err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(body); err != nil {
		slog.Error("error writing response", "error", err)
	}
}
