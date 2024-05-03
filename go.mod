module github.com/crayboi420/chirpy

go 1.22.2

replace github.com/crayboi420/chirpy/internal/database => ./internal/database

require github.com/crayboi420/chirpy/internal/database v0.0.0

require golang.org/x/crypto v0.22.0 // indirect
