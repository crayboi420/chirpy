package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
			ID:    dbUser.ID,
			Email: dbUser.Email,
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
		ID:    user.ID,
		Email: user.Email,
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
			resp = User{ID: user.ID, Email: user.Email}
		}
	}
	if resp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "ID doesn't exist")
	} else {
		respondWithJSON(w, http.StatusOK, resp)
	}
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type inpt struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	type outpt struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := inpt{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Set optional param to 24 hrs
	const SECONDS_IN_DAY int64 = 24 * 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = SECONDS_IN_DAY
	}
	params.ExpiresInSeconds = min(params.ExpiresInSeconds, SECONDS_IN_DAY)

	dbUsers, err := cfg.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
	}
	for _, user := range dbUsers {
		if params.Email == user.Email {
			if bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password)) == nil {
				tkn, err := cfg.JWTToken(user.ID, params.ExpiresInSeconds)
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, "Couldn't generate token: "+err.Error())
					return
				}
				respondWithJSON(w, http.StatusOK, outpt{ID: user.ID, Email: user.Email, Token: tkn})
				return
			} else {
				respondWithError(w, http.StatusUnauthorized, "Password is wrong")
				return
			}
		}
	}
	respondWithError(w, http.StatusNotFound, "user not found")
}

func (cfg *apiConfig) JWTToken(id int, expires_seconds int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			Issuer:    "chirpy",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + expires_seconds,
			Subject:   strconv.Itoa(id),
		},
	)
	tkn, err := token.SignedString([]byte(cfg.jwtSecret))
	return tkn, err
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {

	auth := r.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claims := jwt.StandardClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, 401, "Couldn't parse claims: "+err.Error())
		return
	}
	if !parsed.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token not valid")
		return
	}

	type inpt struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := inpt{}
	decoder.Decode(&params)

	id := claims.Subject
	intID, _ := strconv.Atoi(id)
	cfg.db.UpdateUsers(intID, params.Email, params.Password)

	respondWithJSON(w, http.StatusOK, User{ID: intID, Email: params.Email})
}
