package main

import (
	"net/http"
	"errors"
	"fmt"
	"os"
)

func main(){
	mux:= http.NewServeMux()
	corsMux := middleware(mux)
	serv := http.Server{
		Addr: ":8080",
		Handler: corsMux,
	}

	err:= serv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}