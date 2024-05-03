package main

import (
	"encoding/json"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResp struct {
		Error string `json:"error"`
	}
	resp := errorResp{Error: msg}
	data, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Write(data)
}
