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
	"github.com/TrietFD/student-api/internal/storage/mysql"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := mysql.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	
	// Pass database connection to student handler
	router.HandleFunc("POST /api/students", student.New(storage.Db))
	router.HandleFunc("GET /api/students", student.New(storage.Db))

	// setup server
	server := http.Server{
		Addr:  cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start Server: " + err.Error())
		}
	}()

	<-done
	slog.Info("Shutting down Start Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}