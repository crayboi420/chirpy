module github.com/kanavj/chirpy

go 1.22.2

replace github.com/kanavj/chirpy/internal/database => ./internal/database

require (
	github.com/kanavj/chirpy/internal/database v0.0.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.22.0
)
