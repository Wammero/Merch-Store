package handler

import (
	"net/http"

	"Merch-Store/internal/service"
	"Merch-Store/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	service *service.Service
}

func New(service *service.Service) *API {
	return &API{service: service}
}

func (api *API) SetupRoutes(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Get("/health", healthCheckHandler)

	r.Route("/api", func(router chi.Router) {
		router.Post("/auth", api.Authenticate)

		router.With(jwt.JWTValidator).Group(func(protected chi.Router) {
			protected.Get("/info", api.GetInfo)
			protected.Post("/sendCoin", api.SendCoin)
			protected.Get("/buy/{item}", api.BuyItem)
		})
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
