package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type Compressor struct {
	MinDataLength int
}

// responseRecorder выступает оберткой над http.ResponseWriter.
type gzipResponseWriter struct {
	http.ResponseWriter
	minDataLength int
}

// Write переопределение функции http.ResponseWriter.Write([]byte).
func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	// TODO: переделать длинну порога
	// Сжимаем данные только если их размер больше 75 байт
	if len(data) > w.minDataLength {
		// Получение доступа к пулу
		pool := gzipCompressorPool()
		// Получение writer'а из пула
		i := pool.Get()
		gzip, ok := i.(*gzip.Writer)
		if !ok {
			http.Error(
				w,
				fmt.Sprintf("error getting gzip writer from pool: %s", i.(error)),
				http.StatusInternalServerError,
			)
			return 0, i.(error)
		}
		// Запись сжатых данных в буффер
		var buf bytes.Buffer
		gzip.Reset(&buf)
		_, err := gzip.Write(data)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error writing data to gzip writer: %s", i.(error)),
				http.StatusInternalServerError,
			)
			return 0, err
		}
		// Закрытие и возврат writer'а в пул
		gzip.Close()
		pool.Put(gzip)
		// Установка хедера
		w.Header().Set("Content-Encoding", "gzip")
		slog.Debug(
			"Compressing data",
			slog.Int("raw_data_len", len(data)),
			slog.Int("compressed_data_len", buf.Len()),
		)
		return w.ResponseWriter.Write(buf.Bytes())
	}
	return w.ResponseWriter.Write(data)
}

// GzipCompressor является middleware функцией для использования совместно с chi роутером.
// Сжимает тело ответа, если клиент принимает его в таком виде.
func (c *Compressor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(rw, r)
			return
		}

		grw := &gzipResponseWriter{
			ResponseWriter: rw,
			minDataLength:  c.MinDataLength,
		}

		next.ServeHTTP(grw, r)
	})
}

func gzipCompressorPool() sync.Pool {
	return sync.Pool{
		New: func() any {
			w, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
			if err != nil {
				return err
			}
			return w
		},
	}
}
