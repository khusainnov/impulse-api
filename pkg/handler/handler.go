package handler

import (
	"impulse-api/pkg/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoute() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/signs/{birthday}/{birth_time}/{city}", h.WesternHoroscope).Methods("GET", "POST")
	//r.HandleFunc("/planet-in-houses", h.PlanetsInHouses).Methods("POST")

	return r
}
