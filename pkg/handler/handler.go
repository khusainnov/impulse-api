package handler

import (
	"net/http"

	"impulse-api/pkg/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	//_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	_ "impulse-api/docs"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoute() *mux.Router {
	r := mux.NewRouter()

	docs := r.PathPrefix("/swagger")
	docs.Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	r.HandleFunc("/signs/{birthday}/{birth_time}/{city}/{sex}", h.WesternHoroscope).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/love/{birthday}/{birth_time}/{city}/{sex}", nil).Methods(http.MethodGet, http.MethodPost)
	//r.HandleFunc("/planet-in-houses", h.PlanetsInHouses).Methods("POST")

	return r
}
