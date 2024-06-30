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
	"errors"
	"net/http"
	"strconv"
	"strings"

	serr "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

func (api *API) ordersPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получение имени пользователя
	userID, err := getUserIDFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Чтение тела запроса
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получение строки из тела запроса
	orderStr := buf.String()

	// Проверка номера заказа на валидность
	if !checkLuhn(orderStr) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Преобразование строки в номер заказа
	orderID, err := strconv.ParseUint(orderStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Регистрация заказа
	err = api.orderMgr.Register(r.Context(), userID, orderID)
	switch {
	case errors.Is(err, serr.ErrOrderAlreadyLoaded):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, serr.ErrOrderLoadedByAnotherUsr):
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}
