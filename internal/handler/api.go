package handler

import (
	"net/http"

	"Merch-Store/internal/repository"
	"Merch-Store/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	db *repository.PGRepo
}

func New(db *repository.PGRepo) *API {
	return &API{db: db}
}

func (api *API) SetupRoutes(r *chi.Mux) {
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(router chi.Router) {
		router.Group(func(auth chi.Router) {
			auth.Post("/auth", api.Authenticate)
		})

		router.Group(func(protected chi.Router) {
			protected.Use(jwt.JWTValidator)

			protected.Get("/info", GetInfo)
			protected.Post("/sendCoin", api.SendCoin)
			protected.Get("/buy/{item}", api.BuyItem)
		})
	})
}
