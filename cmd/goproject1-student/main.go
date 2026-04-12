package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rudrapatel51/goproject1-student/internal/config"
	"github.com/rudrapatel51/goproject1-student/internal/http/handlers/student"
	"github.com/rudrapatel51/goproject1-student/internal/storage/postgres"
)

func main() {
	fmt.Println("Welcome to Students API!")

	cfg := config.MustLoad()

	store, err := postgres.New(context.Background(), cfg)
	if err != nil {
		log.Fatal("Failed to connect postgres:", err)
	}
	defer store.Close()

	slog.Info("Connected to Postgres successfully", slog.String("host", cfg.Postgres.Host), slog.Int("port", cfg.Postgres.Port),)

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(store))

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Server Started on", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server:", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server:", err)
	}

	slog.Info("Server Shutdown Successfully")
}
