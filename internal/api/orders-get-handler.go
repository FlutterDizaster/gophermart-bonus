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
	"log/slog"
	"net/http"

	sharederrors "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

func (api *API) ordersGETHandler(w http.ResponseWriter, r *http.Request) {
	// Получение имени пользователя
	userID, err := getUserIDFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получение заказов пользователя
	orders, err := api.orderMgr.Get(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sharederrors.ErrNoOrdersFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal ответа
	body, err := orders.MarshalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(body); err != nil {
		slog.Error("writing response error", "error", err)
	}
}
