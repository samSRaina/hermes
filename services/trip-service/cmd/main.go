package main

import (
	"log"
	"net/http"
	handler "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"

	"github.com/go-chi/chi/v5"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8181")
)

func main() {
	repo := repository.NewInMemRepository()
	service := service.NewTripService(repo)
	handler := handler.HttpHandler{Service: service}

	router := chi.NewRouter()
	router.Post("/preview", handler.HandleTripPreview)

	server := &http.Server{
		Handler: router,
		Addr:    httpAddr,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("http server created: %w", err)
	}

	// to keep the program runnig and preventing docker restarts
	// for {
	// 	time.Sleep(time.Second)
	// }
}
