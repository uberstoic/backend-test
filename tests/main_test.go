package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Serhio/backend-vk-test/internal/config"
	"github.com/Serhio/backend-vk-test/internal/handler"
	"github.com/Serhio/backend-vk-test/internal/service"
	"github.com/Serhio/backend-vk-test/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file for tests: %v", err)
	}
	cfg := config.Load()
	pool, err := postgres.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	postgres.RunMigrations(cfg.DatabaseURL, "../migrations")

	services := service.NewService(pool, cfg.JWTSecret)
	h := handler.NewHandler(services)
	testRouter = h.InitRoutes()

	//run
	code := m.Run()

	//clear
	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE users, ads RESTART IDENTITY")
	if err != nil {
		log.Printf("Failed to truncate tables: %v", err)
	}
	pool.Close()

	os.Exit(code)
}
