package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TrietFD/student-api/internal/config"
	"github.com/TrietFD/student-api/internal/http/handlers/student"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// setup server


	server := http.Server {
		Addr:  	 cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func ()  {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start Server")
		}
	} ()

	<-done

	slog.Info("Shutting down Start Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}