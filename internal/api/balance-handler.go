package api

import (
	"log/slog"
	"net/http"
)

func (api *API) balanceHandler(w http.ResponseWriter, r *http.Request) {
	// Получение имени пользователя
	username, err := getUsernameFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получение баланса пользователя
	balance, err := api.BalanceMgr.Get(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal ответа
	body, err := balance.MarshalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(body); err != nil {
		slog.Error("error writing response", "error", err)
		w.WriteHeader(http.StatusTeapot)
	}
}
