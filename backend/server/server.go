package server

import (
	// "context"
	"log"
	"net/http"
	// "time"
	// "github.com/jackc/pgx/v5"
)

func findByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("recieved request for item: " + id + "\n"))
}

func coolest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting the coolest!\n"))
	/* row := conn.QueryRow(context.Background(),
		"INSERT INTO users (id, username, display_name, role_id, created_at, karma) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		id,
		"Mike",
		"test",
		0,
		time.Now(),
		0,
	)
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("Unable to INSERT: %v\n", err)
		w.WriteHeader(500)
		return
	} */
}

func Setup() {
	router := http.NewServeMux()
	router.HandleFunc("GET /item/{id}", findByID)
	router.HandleFunc("POST /item/{id}", findByID)
	router.HandleFunc("DELETE /item/{id}", findByID)
	router.HandleFunc("GET godotglobes.su/cool", coolest)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	/* stack := middleware.CreateStack(
		middleware.Logging,
		// middleware.AllowCors,
		// middleware.IsAuthed,
		// middleware.CheckPermissions,
	) */

	server := http.Server{
		Addr: ":8080",
		// Handler: stack(router),
		Handler: router,
	}
	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
