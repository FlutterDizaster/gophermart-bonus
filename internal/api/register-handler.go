package api

import (
	"errors"
	"net/http"

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
