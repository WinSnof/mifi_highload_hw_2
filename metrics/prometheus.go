package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		// Используем 'code' для учета ошибок и 'path' для роута
		[]string{"method", "path", "code"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
}

// statusWriter — обертка для http.ResponseWriter, чтобы захватить статус код.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Middleware собирает метрики RPS и latency.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Используем обертку для записи статуса ответа
		sw := &statusWriter{ResponseWriter: w}

		// Передаем запрос дальше в цепочку обработчиков
		next.ServeHTTP(sw, r)

		// Gorilla Mux позволяет получить маршрут, а не просто URL Path
		route := ""
		if currentRoute := mux.CurrentRoute(r); currentRoute != nil {
			route, _ = currentRoute.GetPathTemplate()
		}

		// Если роут не найден (например, 404), используем Path.
		if route == "" {
			route = r.URL.Path
		}

		duration := time.Since(start).Seconds()

		// 1. Учет запросов (RPS и ошибок)
		// Захватываем статус ответа (по умолчанию 200, если не установлен)
		statusCode := sw.status
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		TotalRequests.WithLabelValues(r.Method, route, strconv.Itoa(statusCode)).Inc()

		// 2. Учет latency
		RequestDuration.WithLabelValues(r.Method, route).Observe(duration)
	})
}
