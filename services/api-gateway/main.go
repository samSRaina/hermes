package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hermes/shared/env"

	"github.com/go-chi/chi/v5"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")

	router := chi.NewRouter()
	router.Post("/trip/preview", handleTripPreview)
	router.HandleFunc("/ws/drivers", handleDriverWebSocket)
	router.HandleFunc("/ws/riders", handleRiderWebSocket)

	server := &http.Server{
		Handler: router,
		Addr:    httpAddr,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on port: %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server shutdown due to signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("unsuccessfull graceful shutdown due to %v", err)
			server.Close()
		}

	}
}
