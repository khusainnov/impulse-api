package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"impulse-api/internal/entity"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "impulse-api/docs"
)

const (
	clientID = 618694
)

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

// WesternHoroscope
// Tag DataHandler
// WesternHoroscope - func for signs handler, output type `json`
func (h *Handler) WesternHoroscope(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var API string

	v := mux.Vars(r)
	birthday := v["birthday"]
	birthTime := v["birth_time"]
	city := v["city"]
	sex := v["sex"]

	logrus.Infof("%s // %s // %s ", birthday, birthTime, city)

	// if value haven't birth time then request sends with default time = 12:00
	if len(birthTime) <= 2 {
		API = fmt.Sprintf("%s&date=%s&time=12:00&horo=moon&place=%s", os.Getenv("API"), birthday, city)
	} else {
		API = fmt.Sprintf("%s&date=%s&time=%s&horo=moon&place=%s", os.Getenv("API"), birthday, birthTime, city)
	}

	var dataBody entity.ResponseWithoutTime
	var uprBody entity.ResponseUpr
	var body io.LimitedReader

	client := http.Client{
		Timeout: time.Second * 15,
	}

	// client.Post sends post request on API and gets response body
	resp, err := client.Post(API, "application/json", &body)
	if err != nil {
		logrus.Errorf("Cannot get data from api, due to error: %s", err.Error())
		return
	}

	// h.service.ZodiacApi.GenerateToken generates token with clientID claim
	token, err := h.service.ZodiacApi.GenerateToken(clientID)
	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	w.Header().Set("access-token", token)

	// if we haven't birth time then we go in DataWorkerWithoutTime() and return dataBody as json
	if len(birthTime) <= 2 {
		dataBody, err = h.service.DataWorkerWithoutTime(resp.Body, sex)
		if err != nil {
			_ = json.NewEncoder(w).Encode(&map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(&dataBody)
		return
	} else {
		/*dataBody, err = h.service.DataWorkerWithoutTime(resp.Body, sex)
		if err != nil {
			_ = json.NewEncoder(w).Encode(&map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}*/

		uprBody, err = h.service.DataWorkerWithTime(resp.Body)
		if err != nil {
			_ = json.NewEncoder(w).Encode(&map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(&uprBody)

		/*json.NewEncoder(w).Encode(&map[string]interface{}{
			"without_time": dataBody,
			"with_time":    uprBody,
		})*/

		return
	}
}
