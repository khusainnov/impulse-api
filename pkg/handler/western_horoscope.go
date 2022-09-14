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

const (
	clientID = 618694
)

func (h *Handler) WesternHoroscope(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var API string

	v := mux.Vars(r)
	birthday := v["birthday"]
	birthTime := v["birth_time"]
	city := v["city"]
	sex := v["sex"]

	logrus.Infof("%s // %s // %s ", birthday, birthTime, city)

	if len(birthTime) <= 2 {
		API = fmt.Sprintf("%s&date=%s&time=12:00&horo=moon&place=%s", os.Getenv("API"), birthday, city)
	} else {
		API = fmt.Sprintf("%s&date=%s&time=%s&horo=moon&place=%s", os.Getenv("API"), birthday, birthTime, city)
	}

	var dataBody entity.ResponseWithoutTime
	var uprBody entity.ResponseUpr
	var body io.LimitedReader
	//var fullResp entity.GenResp

	client := http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := client.Post(API, "application/json", &body)
	if err != nil {
		logrus.Errorf("Cannot get data from api, due to error: %s", err.Error())
		return
	}

	token, err := h.service.ZodiacApi.GenerateToken(clientID)
	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	w.Header().Set("access-token", token)

	logrus.Printf("%s \t  %d\n", birthTime, len(birthTime))
	/*if len(birthTime) <= 2 {
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
	}*/

	uprBody, err = h.service.DataWorkerWithTime(resp.Body)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		fmt.Println("WIth time error")
		return
	}

	dataBody, err = h.service.DataWorkerWithoutTime(resp.Body, sex)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		fmt.Println("Without time error")
		return
	}

	json.NewEncoder(w).Encode(&map[string]interface{}{
		"without_time": dataBody,
		"with_time":    uprBody,
	})

	return
}
