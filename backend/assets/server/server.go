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
	router.HandleFunc("POST /assets/", postAsset)
	router.HandleFunc("GET /assets/{id}", getAssetById)
	router.HandleFunc("GET /assets/", getAssets)
	router.HandleFunc("PATCH /assets/", patchAsset)
	router.HandleFunc("DELETE /assets/{id}", deleteAsset)

	router.HandleFunc("GET /docs/assets/", httpSwagger.Handler(
		httpSwagger.URL("/docs/assets/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelRendering":    `"example"`,
			"defaultModelsExpandDepth": "3",
		}),
	))

	server := http.Server{
		Addr:              host,
		ReadHeaderTimeout: 5000 * time.Millisecond,
		ReadTimeout:       5000 * time.Millisecond,
		Handler:           http.TimeoutHandler(router, 5*time.Second, ""),
	}
	log.Printf("Starting server on port: %s\n", host)
	server.ListenAndServe()
}
