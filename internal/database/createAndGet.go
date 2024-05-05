package database

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("user email already exists")
		}
	}
	id := len(dbStructure.Users) + 1

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := User{
		ID:       id,
		Email:    email,
		Password: hash,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Chirps))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) UpdateUsers(targetID int, email string, password string, refresh string) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Find the user by iterating since the key may not match the user ID
	var user *User
	var target int
	for i, u := range dbStructure.Users {
		if u.ID == targetID {
			user = &u
			target = i
			break
		}
	}

	if user == nil {
		return fmt.Errorf("user with ID %d not found", targetID)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Update refresh token logic
	user.RefreshToken, user.RefreshExpiryTime = db.updateRefreshToken(refresh, *user)

	// Always update the email as it's a straightforward string assignment
	user.Email = email

	dbStructure.Users[target] = *user
	// Write updated DB back to storage
	return db.writeDB(dbStructure)
}

// Helper function to handle refresh token logic
func (db *DB) updateRefreshToken(refresh string, user User) (string, time.Time) {
	if refresh != "" {
		// Generate new refresh token with new expiry time
		return refresh, time.Now().Add(24 * 60 * time.Hour) // 60 days
	}
	// If no new token provided, return existing token details
	return user.RefreshToken, user.RefreshExpiryTime
}

func (db *DB) RevokeRefreshToken(refresh string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	for i, usr := range dbStruct.Users {
		if usr.RefreshToken == refresh {
			n_usr := User{Email: usr.Email, ID: usr.ID, Password: usr.Password}
			dbStruct.Users[i] = n_usr
			db.writeDB(dbStruct)
			return nil
		}
	}
	return fmt.Errorf("coulnd't find token %v", refresh)
}
