package database

import (
	"sync"
	"time"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          []byte    `json:"password"`
	RefreshToken      string    `json:"refresh_token"`
	RefreshExpiryTime time.Time `json:"refresh_expiry_time"`
}
