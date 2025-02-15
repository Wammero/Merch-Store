package handler

import (
	"encoding/json"
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

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	writeJSONResponse(w, ErrorResponse{Errors: message}, statusCode)
}

// writeJSONResponse отправляет JSON-ответ с отступами
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, `{
  "errors": "Ошибка при обработке данных"
}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
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
