package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vijay6922/StudentsAPIGoProject/internal/config"
	handlers "github.com/Vijay6922/StudentsAPIGoProject/internal/http/handlers"
)

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", handlers.NewStudent())

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(done)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
		return
	case <-done:
		log.Print("Shutting down server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shut down server gracefully: %v", err)
	}
}
