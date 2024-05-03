package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func Chirp(w http.ResponseWriter, r *http.Request) {
	type inc struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
	}
	dec := json.NewDecoder(r.Body)
	incoming := inc{}
	err := dec.Decode(&incoming)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(incoming.Body) < 140 {
		type validresp struct {
			Cleaned_body string `json:"cleaned_body"`
		}
		resp := validresp{Cleaned_body: cleanWords(incoming.Body)}
		respondWithJSON(w, http.StatusOK, resp)
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}
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

func cleanWords(inc string) string {

	badwords := []string{"kerfuffle", "sharbert", "fornax"}

	wordsinc := strings.Split(inc, " ")
	for i, word := range wordsinc {
		for _, word2 := range badwords {
			if strings.EqualFold(word, word2) {
				wordsinc[i] = "****"
			}
		}
	}
	return strings.Join(wordsinc, " ")
}
