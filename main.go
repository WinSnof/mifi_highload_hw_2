package main

import (
	"fmt"
	"homework_2/handlers"
	"homework_2/metrics"
	"homework_2/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	router := mux.NewRouter()

	router.Use(utils.RateLimitMiddleware)
	router.Use(metrics.Middleware)

	// 3. Регистрация API маршрутов
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	api.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	api.HandleFunc("/users/{id}", handlers.GetUserByIDHandler).Methods("GET")
	api.HandleFunc("/users/{id}", handlers.UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	// 4. Регистрация endpoint для Prometheus (не применяем к нему middleware)
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// 5. Запуск сервера
	port := ":8080"
	fmt.Printf("Сервис запущен на http://localhost%s\n", port)
	fmt.Printf("Метрики доступны на http://localhost%s/metrics\n", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
