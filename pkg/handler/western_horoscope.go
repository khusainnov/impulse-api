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
)

func (h *Handler) WesternHoroscope(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	v := mux.Vars(r)
	birthday := v["birthday"]
	birthTime := v["birth_time"]
	city := v["city"]
	sex := v["sex"]

	logrus.Printf("%s\n%s\n%s\n", birthday, birthTime, city)

	API := fmt.Sprintf("%s&date=%s&time=%s&horo=moon&place=%s", os.Getenv("API"), birthday, birthTime, city)
	//API_CUSPS := fmt.Sprintf("%s&date=%s&time=%s&horo=sun&place=%s",
	//	os.Getenv("API_CUSPS"), birthday, birthTime, city)

	var dataBody entity.Summary
	var uprBody entity.ResponseUpr
	var body io.LimitedReader

	client := http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := client.Post(API, "application/json", &body)
	if err != nil {
		logrus.Errorf("Cannot get data from api, due to error: %s", err.Error())
		return
	}

	//resp_cusp, err := client.Post(API_CUSPS, "application/json", &body)
	//if err != nil {
	//	logrus.Errorf("Cannot get data from cusps_api, due to error: %s", err.Error())
	//	return
	//}

	token, err := h.service.ZodiacApi.GenerateToken(618694)
	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	w.Header().Set("access-token", token)

	logrus.Printf("%s \t  %d\n", birthTime, len(birthTime))
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

		uprBody, err = h.service.DataWorkerWithTime(resp.Body)
		if err != nil {
			_ = json.NewEncoder(w).Encode(&map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(&uprBody)
		return
	}
}
