package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	logrus.Printf("%s\n%s\n%s\n", birthday, birthTime, city)

	API := fmt.Sprintf("http://astro2022.fun/hooksp.php?action=western_horoscope&date=%s&time=%s&horo=moon&place=%s", birthday, birthTime, city)

	var dataBody entity.Summary
	var body io.LimitedReader

	//api := os.Getenv("API")

	client := http.Client{
		Timeout: time.Second * 15,
	}

	resp, err := client.Post(API, "application/json", &body)
	if err != nil {
		logrus.Errorf("Cannot get data from api, due to error: %s", err.Error())
		return
	}

	token, err := h.service.ZodiacApi.GenerateToken(618694, "6b03d1dfcbb5d09e704badf9730e3fca")
	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	w.Header().Set("access-token", token)

	dataBody, err = h.service.DataWorker(resp.Body)

	if err != nil {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	/*dBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Cannot read body, due to error: %s", err.Error())
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		logrus.Errorf("Cannot unmarshall body, due to error: %s", err.Error())
	}*/

	json.NewEncoder(w).Encode(&dataBody)

	return
}
