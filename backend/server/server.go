package server

import (
	"log"
	"net/http"
)

func coolest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting the coolest!\n"))
}

func Setup() {
	router := http.NewServeMux()
	router.HandleFunc("POST /demos/", postDemo)
	router.HandleFunc("GET /demos/{id}", getDemoById)
	// router.HandleFunc("GET /demos", GetDemos)
	router.HandleFunc("DELETE /demos/{id}", coolest)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
