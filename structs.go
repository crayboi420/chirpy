package main

import "github.com/crayboi420/chirpy/internal/database"

type apiConfig struct {
	fileserverHits int
	db             database.DB
	jwtSecret      string
	polkaKey       string
}

type Chirp struct {
	AuthorID int    `json:"author_id"`
	ID       int    `json:"id"`
	Body     string `json:"body"`
}

type User struct {
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
	Email       string `json:"email"`
	// Password string `json:"password"`
}
