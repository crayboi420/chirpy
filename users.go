package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerUsersRetrieve(w http.ResponseWriter, _ *http.Request) {

	dbUsers, err := cfg.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
		return
	}

	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:          dbUser.ID,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
			// Password: dbUser.Password,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
		// Password: user.Password,
	})
}

func (cfg *apiConfig) handlerUserRetrieve(w http.ResponseWriter, r *http.Request) {
	req := r.PathValue("userID")

	reqInt, err := strconv.Atoi(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert number to int")
	}

	dbUsers, err := cfg.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
		return
	}

	resp := User{}
	for _, user := range dbUsers {
		if user.ID == reqInt {
			resp = User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
		}
	}
	if resp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "ID doesn't exist")
	} else {
		respondWithJSON(w, http.StatusOK, resp)
	}
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	
	authstr := strings.TrimPrefix(r.Header.Get("Authorization"),"ApiKey ")
	if authstr!= cfg.polkaKey{
		respondWithError(w,http.StatusUnauthorized,"wrong api key")
		return
	}
	
	type inpt struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	n_inpt := inpt{}
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&n_inpt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request :"+err.Error())
		return
	}
	event := n_inpt.Event
	if event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "{}")
		return
	}

	id := n_inpt.Data.UserID
	err = cfg.db.UpdateRed(id, true)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	} else {
		respondWithJSON(w, http.StatusOK, "{}")
	}

}
