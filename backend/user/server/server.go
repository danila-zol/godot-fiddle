package server

import (
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"

	"github.com/swaggo/http-swagger"
)

func Setup(host string) {
	router := http.NewServeMux()
	router.HandleFunc("POST /user/", postUser)
	router.HandleFunc("GET /user/{id}", getUserById)
	router.HandleFunc("GET /user/", getUsers)
	router.HandleFunc("PATCH /user/", patchUser)
	router.HandleFunc("DELETE /user/{id}", deleteUser)

	router.HandleFunc("POST /roles/", postRole)
	router.HandleFunc("GET /roles/{id}", getRoleById)
	router.HandleFunc("GET /roles/", getRoles)
	router.HandleFunc("PATCH /roles/", patchRole)
	router.HandleFunc("DELETE /roles/{id}", deleteRole)

	router.HandleFunc("GET /docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelRendering":    `"example"`,
			"defaultModelsExpandDepth": "3",
		}),
	))

	server := http.Server{
		Addr:              host,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		Handler:           http.TimeoutHandler(router, time.Second, ""),
	}
	log.Printf("Starting server on port: %s\n", host)
	server.ListenAndServe()
}
