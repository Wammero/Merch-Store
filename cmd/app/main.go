package main

import (
	"fmt"
	"log"
	"net/http"

	"Merch-Store/cmd/migrate"
	"Merch-Store/internal/handler"
	"Merch-Store/internal/repository"
	"Merch-Store/internal/service"
	"Merch-Store/pkg/configload"
	"Merch-Store/pkg/jwt"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := configload.LoadConfig("/app/configs/config.yaml")
	connStr := buildDBConnectionString(cfg)

	repo, err := repository.New(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()
	log.Println("Successfully connected to the database!")

	jwt.SetSecret(cfg.JWT.Secret)
	applyMigrations(connStr)

	svc := service.New(repo)
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.URLFormat)
	handler.New(svc).SetupRoutes(r)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func buildDBConnectionString(cfg *configload.Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)
}

func applyMigrations(connStr string) {
	m, err := migrate.CallMigrations(connStr)
	if err != nil {
		log.Fatalf("Error initializing migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("Migrations are already applied")
		} else {
			log.Fatalf("Error applying migrations: %v", err)
		}
	} else {
		log.Println("Migrations successfully applied")
	}
}
