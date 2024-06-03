package api

import (
	"log/slog"
	"net/http"
)

func (api *API) ordersGETHandler(w http.ResponseWriter, r *http.Request) {
	// Получение имени пользователя
	username, err := getUsernameFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Получение заказов пользователя
	orders, err := api.orderMgr.Get(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
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
		w.WriteHeader(http.StatusTeapot)
	}
}
