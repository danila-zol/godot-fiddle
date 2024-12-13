package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

//	@Summary	Creates a new demo.
//	@Tags		Demos
//	@Accept		application/json
//	@Produce	application/json
//	@Param		Demo	body		database.Demo	true	"Create Demo"
//	@Success	200		{object}	ResponseHTTP{data=database.Demo}
//	@Failure	400		{object}	ResponseHTTP{}
//	@Failure	500		{object}	ResponseHTTP{}
//	@Router		/demos/ [post]
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

	newDemo, err := operations.CreateDemo(demo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateDemo operation\n"))
		log.Printf("Error in CreateDemo operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newDemo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in postDemo operation\n"))
		log.Printf("Error in postDemo operation \n%s", err)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte("Demo successfully created under ID " + newDemo.ID + "!\n"))
}

//	@Summary	Fetches a demo by its ID.
//	@Tags		Demos
//	@Accept		text/plain
//	@Produce	application/json
//	@Param		id	path		string	true	"Get Demo of ID"
//	@Success	200	{object}	ResponseHTTP{data=database.Demo}
//	@Failure	400	{object}	ResponseHTTP{}
//	@Failure	500	{object}	ResponseHTTP{}
//	@Router		/demos/{id} [get]
func getDemoById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	demo, err := operations.FindFirstDemo(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Demo not found!\n"))
		log.Printf("Error: Demo not found!\n%s", err)
		return
	}
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
		w.Write([]byte("Error in getDemoById operation\n"))
		log.Printf("Error in getDemoById operation \n%s", err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("recieved request for item: " + id + "\n"))
}

//	@Summary	Fetches all demos.
//	@Tags		Demos
//	@Produce	application/json
//	@Success	200	{object}	ResponseHTTP{data=[]database.Demo}
//	@Failure	400	{object}	ResponseHTTP{}
//	@Failure	500	{object}	ResponseHTTP{}
//	@Router		/demos/ [get]
func getDemos(w http.ResponseWriter, r *http.Request) {
	demo, err := operations.FindDemos()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Demos not found!\n"))
		log.Printf("Error: Demos not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindDemos operation\n"))
		log.Printf("Error in FindDemos operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(demo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getDemos operation\n"))
		log.Printf("Error in getDemos operation \n%s", err)
		return
	}
	w.WriteHeader(200)
}

//	@Summary	Updates a demo.
//	@Tags		Demos
//	@Accept		application/json
//	@Produce	application/json
//	@Param		Demo	body		database.Demo	true	"Update Demo"
//	@Success	200		{object}	ResponseHTTP{data=database.Demo}
//	@Failure	400		{object}	ResponseHTTP{}
//	@Failure	500		{object}	ResponseHTTP{}
//	@Router		/demos/ [patch]
func patchDemo(w http.ResponseWriter, r *http.Request) {
	var demo operations.Demo
	err := json.NewDecoder(r.Body).Decode(&demo)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchDemo handler\n"))
		log.Printf("Error in patchDemo handler \n%s", err)
		return
	}
	demo.Updated_at = time.Now()

	updDemo, err := operations.UpdateDemo(demo)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Demos not found!\n"))
		log.Printf("Error: Demos not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateDemo operation\n"))
		log.Printf("Error in UpdateDemo operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updDemo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in patchDemo handler\n"))
		log.Printf("Error in patchDemo handler \n%s", err)
		return
	}
	w.WriteHeader(200)
}

//	@Summary	Deletes the specified demo.
//	@Tags		Demos
//	@Accept		text/plain
//	@Produce	text/plain
//	@Param		id	path		string	true	"Delete Demo of ID"
//	@Success	200	{object}	ResponseHTTP{}
//	@Failure	400	{object}	ResponseHTTP{}
//	@Failure	500	{object}	ResponseHTTP{}
//	@Router		/demos/{id} [delete]
func deleteDemo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteDemo(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Demo not found!\n"))
		log.Printf("Error: Demo not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstDemo operation\n"))
		log.Printf("Error in FindFirstDemo operation \n%s", err)
		return
	}

	w.WriteHeader(200)
}
