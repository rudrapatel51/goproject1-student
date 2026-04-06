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
)

func main() {
	fmt.Println("Welcome to Students API!")

	cfg := config.MustLoad()

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address, // ✅ correct field
		Handler: router,
	}

	slog.Info("Server Started on", slog.String("address", cfg.HTTPServer.Address)) // ✅ fixed

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed { // ✅ fixed
			log.Fatal("Error starting server:", err)
		}
	}()

	<-done // ✅ block until shutdown signal received

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server:", err)
	}

	slog.Info("Server Shutdown Successfully")
}
