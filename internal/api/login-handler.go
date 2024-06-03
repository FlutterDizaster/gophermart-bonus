package api

import (
	"errors"
	"net/http"

	sharederrors "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
)

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных для авторизации пользователя
	cred, err := getCredentialsFromReq(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Авторизация пользователя
	token, err := api.userMgr.Login(r.Context(), cred)
	if err != nil {
		if errors.Is(err, sharederrors.ErrWrongLoginOrPassword) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// TODO: Добавить время жизни cookie
	// Создание cookie
	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
	}

	// Установка cookie
	http.SetCookie(w, cookie)

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}
