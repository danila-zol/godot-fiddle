package server

import (
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"

	"github.com/swaggo/http-swagger"
)

func Setup() {
	router := http.NewServeMux()
	router.HandleFunc("POST /demos/", postDemo)
	router.HandleFunc("GET /demos/{id}", getDemoById)
	router.HandleFunc("GET /demos/", getDemos)
	router.HandleFunc("PATCH /demos/", patchDemo)
	router.HandleFunc("DELETE /demos/{id}", deleteDemo)

	router.HandleFunc("GET /docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelRendering":    `"example"`,
			"defaultModelsExpandDepth": "3",
		}),
	))

	server := http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		Handler:           http.TimeoutHandler(router, time.Second, ""),
	}
	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
