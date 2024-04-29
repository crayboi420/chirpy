package main

import "net/http"

func middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Methods","GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers","*")
		if r.Method == "OPTIONS"{
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w,r)
	})
}

func createServer(addr string) http.Server {

	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir("./files/")))
	corsMux := middleware(mux)
	return http.Server{
		Addr: addr,
		Handler: corsMux,
	}
}