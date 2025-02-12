package main

import (
	"fmt"
	"log"
	"net/http"

	"Merch-Store/cmd/migrate"
	"Merch-Store/internal/handler"
	repo "Merch-Store/internal/repository"
	config "Merch-Store/pkg/configload"
	"Merch-Store/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig("/app/configs/config.yaml")
	connstr := buildDBConnectionString(cfg)
	repo := connectToDatabase(connstr)
	defer repo.Close()

	jwt.SetSecret(cfg.JWT.Secret)

	applyMigrations(connstr)

	api := handler.New(repo)

	r := chi.NewRouter()
	api.SetupRoutes(r)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func buildDBConnectionString(cfg *config.Config) string {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)
	return connStr
}

func connectToDatabase(connstr string) *repo.PGRepo {
	repo, err := repo.New(connstr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database!")
	return repo
}

func applyMigrations(connstr string) {
	m, err := migrate.CallMigrations(connstr)
	if err != nil {
		log.Fatalf("Ошибка при создании мигратора: %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("Миграции уже применены")
		} else {
			log.Fatalf("Ошибка при применении миграции: %v", err)
		}
	} else {
		log.Println("Миграции успешно применены")
	}
}
