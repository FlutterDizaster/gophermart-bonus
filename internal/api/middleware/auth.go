// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License.
//
// This file uses a third-party package jwt licensed under MIT License.
//
// See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package middleware

import (
	"context"
	"net/http"

	ctxkeys "github.com/FlutterDizaster/gophermart-bonus/internal/context-keys"
	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
)

type TokenResolver interface {
	DecryptToken(tokenString string) (*models.Claims, error)
}

type AuthMiddleware struct {
	Resolver    TokenResolver
	PublicPaths []string
}

var _ Middleware = &AuthMiddleware{}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка URL на наличие в списке публичных
		for i := range m.PublicPaths {
			// Если есть совпадение, то пропуск проверки
			if m.PublicPaths[i] == r.URL.Path {
				next.ServeHTTP(w, r)
			}
		}

		// Получение cookie
		cookies := r.Cookies()
		for i := range cookies {
			if cookies[i].Name == "Auth" {
				// Проверка куки
				if userID, ok := m.checkCookie(cookies[i]); ok {
					// Сохранение userID в контекст
					ctx := context.WithValue(r.Context(), ctxkeys.UserID, userID)

					// Создание нового запроса с контекстом
					req := r.WithContext(ctx)

					// Передача управления следующему хендлеру с новым запросом
					next.ServeHTTP(w, req)
					return
				}
			}
		}
		// Отправка ответа
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// Проверяет cookie на валидность и возвращает первым аргументом userID,
// а вторым статус проверки.
// Если проверка пройдена успешно, то userID будет содержать имя пользователя,
// а статус будет true. Иначе userID будет пустой строкой, а статус false.
func (m *AuthMiddleware) checkCookie(cookie *http.Cookie) (uint64, bool) {
	// Проверка cookie на валидность
	err := cookie.Valid()
	if err != nil {
		return 0, false
	}

	// Расшифровка токена
	claims, err := m.Resolver.DecryptToken(cookie.Value)
	if err != nil {
		return 0, false
	}

	return claims.UserID, true
}
