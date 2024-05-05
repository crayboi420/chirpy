package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type inpt struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	type outpt struct {
		ID      int    `json:"id"`
		Email   string `json:"email"`
		Token   string `json:"token"`
		Refresh string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := inpt{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Set optional param to 24 hrs
	const expires_in_seconds int64 = 60 * 60

	dbUsers, err := cfg.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
	}
	for _, user := range dbUsers {
		if params.Email == user.Email {
			if bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password)) == nil {
				tkn, err := cfg.JWTToken(user.ID, expires_in_seconds)
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, "Couldn't generate token: "+err.Error())
					return
				}
				refr := generateRefresh(32)
				cfg.db.UpdateUsers(user.ID, user.Email, params.Password, refr)
				respondWithJSON(w, http.StatusOK, outpt{ID: user.ID, Email: user.Email, Token: tkn, Refresh: refr})
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

	claims,err := cfg.checkHeader(r)

	if err!=nil{
		respondWithError(w,http.StatusUnauthorized,err.Error())
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
	cfg.db.UpdateUsers(intID, params.Email, params.Password, "")

	respondWithJSON(w, http.StatusOK, User{ID: intID, Email: params.Email})
}

func generateRefresh(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type outpt struct {
		Token string `json:"token"`
	}

	auth := r.Header.Get("Authorization")
	refresh := strings.TrimPrefix(auth, "Bearer ")
	if len(refresh) == 0 {
		respondWithError(w, http.StatusBadRequest, "No token provided")
		return
	}

	dbUsers, err := cfg.db.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}

	for _, usr := range dbUsers {
		if refresh == usr.RefreshToken {
			if time.Now().After(usr.RefreshExpiryTime) {
				respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
				return
			} else {
				tkn, err := cfg.JWTToken(usr.ID, 60*60*24)
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, "Couldn't Generate token")
					return
				}
				respondWithJSON(w, http.StatusOK, outpt{Token: tkn})
				return
			}
		}
	}
	respondWithError(w, http.StatusUnauthorized, "Token doesn't exist")
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	refresh := strings.TrimPrefix(auth, "Bearer ")

	if len(refresh) == 0 {
		respondWithError(w, http.StatusBadRequest, "No token provided")
		return
	}
	err := cfg.db.RevokeRefreshToken(refresh)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}

func (cfg *apiConfig) checkHeader(r *http.Request) (*jwt.StandardClaims,error){
	
	
	auth := r.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	claims := jwt.StandardClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		return &claims,fmt.Errorf("coulnd't parse claims")
	}
	if !parsed.Valid {
		return &claims,fmt.Errorf("not a valid token")
	}
	return &claims,nil
}