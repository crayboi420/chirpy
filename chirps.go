package main

import (
	"strings"
	"net/http"
	"sort"
	"encoding/json"
	"strconv"
	"errors"
)

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

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, _ *http.Request) {

	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	req := r.PathValue("chirpID")

	reqInt, err := strconv.Atoi(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert number to int")
	}

	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	resp := Chirp{}
	for _, chirp := range dbChirps {
		if chirp.ID == reqInt {
			resp = Chirp{ID: chirp.ID, Body: chirp.Body}
		}
	}
	if resp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "ID doesn't exist")
	} else {
		respondWithJSON(w, http.StatusOK, resp)
	}
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned := cleanWords(body)
	return cleaned, nil
}