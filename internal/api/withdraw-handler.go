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
	"strconv"

	serr "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

func (api *API) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	// Получение запроса на списание
	withdraw, err := getWithdrawFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	withdraw.OrderID, err = strconv.ParseUint(withdraw.StringOrderID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получение имени пользователя
	userID, err := getUserIDFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Debug(
		"new withdraw request",
		slog.Group(
			"details",
			slog.Uint64("user id", userID),
			slog.Group(
				"withdraw",
				slog.Uint64("orderID", withdraw.OrderID),
				slog.Float64("sum", withdraw.Sum),
			),
		),
	)

	// Проведение списания
	err = api.BalanceMgr.ProcessWithdraw(r.Context(), userID, withdraw)
	switch {
	case errors.Is(err, serr.ErrNotEnoughFunds):
		w.WriteHeader(http.StatusPaymentRequired)
	case errors.Is(err, serr.ErrWithdrawNotAllowed):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, serr.ErrWrongOrderID):
		w.WriteHeader(http.StatusUnprocessableEntity)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
