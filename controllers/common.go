package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"rouletteapi/prometheus"
)

func writeJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		panic(err)
	}
}

func writeErrorWithMsg(w http.ResponseWriter, r *http.Request, err error) {
	statusCode := http.StatusInternalServerError
	prometheus.ErrorCounter.Inc()

	http.Error(w, err.Error(), statusCode)
}

func GenerateHash() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	return s
}