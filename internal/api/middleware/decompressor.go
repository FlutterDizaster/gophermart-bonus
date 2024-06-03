package middleware

import (
	"compress/gzip"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Decompressor struct{}

// GzipUncompressor является middleware функцией для использования совместно с chi роутером.
// Распаковывает тело запроса, если клиент отправил его в таком виде.
func (d *Decompressor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") || r.Body == nil {
			next.ServeHTTP(rw, r)
			return
		}

		// Создание ридера
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			slog.Error("error creating gzip reader", "error", err)
			http.Error(
				rw,
				fmt.Sprintf("error creating gzip reader: %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Подмена body
		r.Body = reader

		next.ServeHTTP(rw, r)
	})
}
