package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func main() {

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            false,
		MaxAge:           300,
	})

	router.Use(cors.Handler)

	//Handle Graceful Shutdown using server.Shutdown
	server := http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	//db := connectDatabase(config.GetConfig())
	//router.Mount("/api/locations", routes.Routes(db))

	idleConnsClosed := make(chan struct{})
	go GracefulShutdownWatcher(idleConnsClosed, &server)

	log.Printf("Serving at %s \n", "4000")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe error: %v", err)
	}

	<-idleConnsClosed
}

func GracefulShutdownWatcher(idleConnsClosed chan struct{}, server *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	log.Println("Shutting down. Goodbye..")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}
	close(idleConnsClosed)
}
