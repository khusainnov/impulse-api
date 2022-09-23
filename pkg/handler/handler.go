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

// @Summary 	WesternHoroscope
// @Tags 		data handler
// @Description data handler gets data from vars in url and then request to astrobot API for getting all data about planets in this date
// @Accept  	json
// @Produce 	json
// @Param   	dataBody entity.ResponseWithoutTime
// @Param   	uprBody entity.ResponseUpr
// @Success 	200 {object} entity.ResponseWithoutTime
// @Success 	200 {object} entity.ResponseUpr
// @Failure 	500 {object} json.Encode()
// @Router  	/signs/birthday/birth_time/city/sex [get]

func (h *Handler) InitRoute() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/signs/{birthday}/{birth_time}/{city}/{sex}", h.WesternHoroscope).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/zodiac", h.Zodiac).Methods(http.MethodGet, http.MethodPost)

	docs := r.PathPrefix("/swagger")
	docs.Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	return r
}
