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
	router.HandleFunc("POST /topics/", postTopic)
	router.HandleFunc("GET /topics/{id}", getTopicById)
	router.HandleFunc("GET /topics/", getTopics)
	router.HandleFunc("PATCH /topics/", patchTopic)
	router.HandleFunc("DELETE /topics/{id}", deleteTopic)

	router.HandleFunc("POST /threads/", postThread)
	router.HandleFunc("GET /threads/{id}", getThreadById)
	router.HandleFunc("GET /threads/", getThreads)
	router.HandleFunc("PATCH /threads/", patchThread)
	router.HandleFunc("DELETE /threads/{id}", deleteThread)

	router.HandleFunc("POST /messages/", postMessage)
	router.HandleFunc("GET /messages/{id}", getMessageById)
	router.HandleFunc("GET /messages/", getMessages)
	router.HandleFunc("PATCH /messages/", patchMessage)
	router.HandleFunc("DELETE /messages/{id}", deleteMessage)

	router.HandleFunc("GET /docs/forum/", httpSwagger.Handler(
		httpSwagger.URL("/docs/forum/doc.json"),
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
