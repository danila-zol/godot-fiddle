package server

import (
	"log"
	"net/http"

	_ "game-hangar/docs"
	"github.com/swaggo/http-swagger"
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

	router.HandleFunc("GET /docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelRendering":    `"example"`,
			"defaultModelsExpandDepth": "3",
		}),
	))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
