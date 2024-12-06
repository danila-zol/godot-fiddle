package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"
	"time"
	// "github.com/google/uuid"
)

func postDemo(w http.ResponseWriter, r *http.Request) {
	var demo operations.Demo
	err := json.NewDecoder(r.Body).Decode(&demo)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostDemo handler\n"))
		log.Printf("Error in PostDemo handler \n%s", err)
		return
	}

	// demo.ID = "demo_" + uuid.NewString()
	demo.Created_at, demo.Updated_at = time.Now(), time.Now()
	demo.Upvotes, demo.Downvotes = 0, 0

	id, err := operations.CreateDemo(demo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateDemo operation\n"))
		log.Printf("Error in CreateDemo operation \n%s", err)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte("Demo successfully created under ID " + id + "!\n"))
}

func getDemoById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	demo, err := operations.FindFirstDemo(id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstDemo operation\n"))
		log.Printf("Error in FindFirstDemo operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(demo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in GetDemoById operation\n"))
		log.Printf("Error in GetDemoById operation \n%s", err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("recieved request for item: " + id + "\n"))
}
