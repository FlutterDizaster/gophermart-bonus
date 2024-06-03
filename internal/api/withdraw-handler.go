package api

import (
	"errors"
	"net/http"

	serr "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

func (api *API) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	// Получение запроса на списание
	withdraw, err := getWithdrawFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получение имени пользователя
	username, err := getUsernameFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Проведение списания
	err = api.BalanceMgr.ProcessWithdraw(r.Context(), username, withdraw)
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
