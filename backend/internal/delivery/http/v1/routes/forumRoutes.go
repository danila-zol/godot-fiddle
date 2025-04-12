package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	"github.com/labstack/echo/v4"
)

type ForumRoutes struct {
	handler *handlers.ForumHandler
}

func NewForumRoutes(h *handlers.ForumHandler) *ForumRoutes {
	return &ForumRoutes{
		handler: h,
	}
}

func (r *ForumRoutes) InitRoutes(e *echo.Echo) {
	topicGroup := e.Group("/game-hangar/v1/topics")

	protectedTopicGroup := topicGroup.Group("")

	protectedTopicGroup.POST("", r.handler.PostTopic)
	topicGroup.GET("/:id", r.handler.GetTopicById)
	topicGroup.GET("", r.handler.GetTopics)
	protectedTopicGroup.PATCH("/:id", r.handler.PatchTopic)
	protectedTopicGroup.DELETE("/:id", r.handler.DeleteTopic)

	threadGroup := e.Group("/game-hangar/v1/threads")

	protectedThreadGroup := threadGroup.Group("")

	protectedThreadGroup.POST("", r.handler.PostThread)
	threadGroup.GET("/:id", r.handler.GetThreadById)
	threadGroup.GET("", r.handler.GetThreads)
	protectedThreadGroup.PATCH("/:id", r.handler.PatchThread)
	protectedThreadGroup.DELETE("/:id", r.handler.DeleteThread)

	messageGroup := e.Group("/game-hangar/v1/messages")

	protectedMessageGroup := messageGroup.Group("")

	protectedMessageGroup.POST("", r.handler.PostMessage)
	messageGroup.GET("/:id", r.handler.GetMessageById)
	messageGroup.GET("", r.handler.GetMessages)
	protectedMessageGroup.PATCH("/:id", r.handler.PatchMessage)
	protectedMessageGroup.DELETE("/:id", r.handler.DeleteMessage)
}
