package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"impulse-api/internal/entity"

	"github.com/sirupsen/logrus"
)

func (h *Handler) PlanetsInHouses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dataBody entity.Summary
	var body io.LimitedReader

	api := os.Getenv("API")

	client := http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := client.Post(api, "application/json", &body)
	if err != nil {
		logrus.Errorf("Cannot get data from api, due to error: %s", err.Error())
		return
	}

	dataBody, err = h.service.DataWorker(resp.Body)
	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&map[string]interface{}{
		"planet": dataBody.Planets[0].Name,
		"house":  dataBody.Houses[0].House,
		"sign":   dataBody.Houses[0].Sign,
	})
}
