package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Serhio/backend-vk-test/internal/config"
	"github.com/Serhio/backend-vk-test/internal/handler"
	"github.com/Serhio/backend-vk-test/internal/service"
	"github.com/Serhio/backend-vk-test/internal/storage/postgres"
)

func main() {
	// init config
	cfg := config.Load()

	// init db
	pool, err := postgres.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// run migrations
	postgres.RunMigrations(cfg.DatabaseURL, "migrations")

	// init services
	services := service.NewService(pool, cfg.JWTSecret)

	// init handlers
	h := handler.NewHandler(services)

	// run server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: h.InitRoutes(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)
	log.Println("Marketplace API starting...")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
