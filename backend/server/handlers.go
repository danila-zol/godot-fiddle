package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"
	"github.com/google/uuid"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

//	@Summary	Creates a new demo.
//	@Tags		Demos
//	@Accept		json
//	@Produce	json
//	@Param		Demo	body		database.Demo	true	"Create Demo"
//	@Success	200		{object}	ResponseHTTP{data=database.Demo}
//	@Failure	400		{object}	ResponseHTTP{}
//	@Failure	500		{object}	ResponseHTTP{}
//	@Router		/demos [post]
func postDemo(w http.ResponseWriter, r *http.Request) {
	var demo operations.Demo
	err := json.NewDecoder(r.Body).Decode(&demo)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostDemo handler\n"))
		log.Printf("Error in PostDemo handler \n%s", err)
		return
	}

	demo.ID = "demo_" + uuid.NewString()
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

//	@Summary	Fetches a demo by its ID.
//	@Tags		Demos
//	@Accept		text/plain
//	@Produce	json
//	@Param		id	path		string	true	"Get Demo by ID"
//	@Success	200	{object}	ResponseHTTP{data=database.Demo}
//	@Failure	400	{object}	ResponseHTTP{}
//	@Failure	500	{object}	ResponseHTTP{}
//	@Router		/demos/{id} [get]
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
