package handler

import "net/http"

func (h *Handler) Zodiac(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Zodiac"))
}
