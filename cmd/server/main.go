package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	usecase "github.com/dementievme/pull-request-service/internal/application/usecase"
	domain "github.com/dementievme/pull-request-service/internal/domain/service"
	postgres "github.com/dementievme/pull-request-service/internal/infrastructure/db/postgresql"
	api "github.com/dementievme/pull-request-service/internal/infrastructure/http/api/v1"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Получение переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open postgres connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	repos := postgres.NewRepositories(db)

	services := domain.NewServices(repos.PullRequestRepo, repos.TeamRepo, repos.UserRepo)

	useCases := usecase.NewUseCases(services.PullRequestService, services.TeamService, services.UserService)

	r := gin.Default()
	api.RegisterRoutes(r, useCases)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
