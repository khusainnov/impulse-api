package handler

import (
	"encoding/json"
	"net/http"

	_ "impulse-api/docs"
)

// @Summary 	Zodiac
// @Tags 		zodiac handler
// @Description unused function
// @Accept  	json
// @Produce 	json
// @Success 	200 {object} w.Write()
// @Failure 	500 {object} json.Encode()
// @Router  	/zodiac [get, post]

func (h *Handler) Zodiac(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("Zodiac"))
	if err != nil {
		_ = json.NewEncoder(w).Encode(&map[string]interface{}{
			"code": http.StatusInternalServerError,
		})
		return
	}

	return
}
