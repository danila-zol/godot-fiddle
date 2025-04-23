package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"strconv"
	"time"

	_ "gamehangar/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ForumHandler struct {
	logger     echo.Logger
	repository ForumRepository
	validator  *validator.Validate
}

func NewForumHandler(e *echo.Echo, repo ForumRepository, v *validator.Validate) *ForumHandler {
	return &ForumHandler{
		logger:     e.Logger,
		repository: repo,
		validator:  v,
	}
}

// @Summary	Creates a new topic.
// @Tags		Topics
// @Accept		application/json
// @Produce	application/json
// @Param		Topic	body		models.Topic	true	"Create Topic"
// @Success	200		{object}	ResponseHTTP{data=models.Topic}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/topics [post]
func (h *ForumHandler) PostTopic(c echo.Context) error {
	var topic models.Topic

	err := c.Bind(&topic)
	if err != nil {
		h.logger.Printf("Error in PostTopic handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostTopic handler")
	}

	err = h.validator.Struct(&topic)
	if err != nil {
		h.logger.Printf("Error in PostTopic handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PostTopic handler")
	}

	newTopic, err := h.repository.CreateTopic(topic)
	if err != nil {
		h.logger.Printf("Error in CreateTopic repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateTopic repository")
	}

	return c.JSON(http.StatusOK, &newTopic)
}

// @Summary	Fetches a topic by its ID.
// @Tags		Topics
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Topic of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Topic}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/topics/{id} [get]
func (h *ForumHandler) GetTopicByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetTopicByID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetTopicByID handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	topic, err := h.repository.FindTopicByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Topic not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Topic not found!")
		}
		h.logger.Printf("Error in FindTopicByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindTopicByID repository")
	}

	return c.JSON(http.StatusOK, &topic)
}

// @Summary	Fetches all topics.
// @Tags		Topics
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Topic}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/topics [get]
func (h *ForumHandler) GetTopics(c echo.Context) error {
	topics, err := h.repository.FindTopics()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Topic not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Topic not found!")
		}
		h.logger.Printf("Error in FindFirstTopic repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindTopics repository")
	}

	return c.JSON(http.StatusOK, &topics)
}

// @Summary	Updates an topic.
// @Tags		Topics
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Topic of ID"
// @Param		Topic	body		models.Topic	true	"Update Topic"
// @Success	200		{object}	ResponseHTTP{data=models.Topic}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/topics/{id} [patch]
func (h *ForumHandler) PatchTopic(c echo.Context) error {
	var topic models.Topic
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in PatchTopic handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchTopic handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	if err := c.Bind(&topic); err != nil {
		h.logger.Printf("Error in PatchTopic handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchTopic handler")
	}

	err = h.validator.Struct(&topic)
	if err != nil {
		h.logger.Printf("Error in PatchTopic handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchTopic handler")
	}

	updTopic, err := h.repository.UpdateTopic(int(id), topic)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Topic not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Topic not found!")
		}
		h.logger.Printf("Error in UpdateTopic repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateTopic repository")
	}

	return c.JSON(http.StatusOK, &updTopic)
}

// @Summary	Deletes the specified topic.
// @Tags		Topics
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Topic of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/topics/{id} [delete]
func (h *ForumHandler) DeleteTopic(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in DeleteTopic handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in DeleteTopic handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteTopic(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Topic not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Topic not found!")
		}
		h.logger.Printf("Error in DeleteTopic repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteTopic repository")
	}

	return c.String(http.StatusOK, "Topic sucessfully deleted!")
}

// @Summary	Creates a new thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		Thread	body		models.Thread	true	"Create Thread"
// @Success	200		{object}	ResponseHTTP{data=models.Thread}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/threads [post]
func (h *ForumHandler) PostThread(c echo.Context) error {
	var thread models.Thread

	err := c.Bind(&thread)
	if err != nil {
		h.logger.Printf("Error in PostThread handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostThread handler")
	}
	thread.Method = "POST"

	if thread.CreatedAt == nil || thread.UpdatedAt == nil {
		currentTime := time.Now()
		thread.CreatedAt, thread.UpdatedAt = &currentTime, &currentTime
	}
	if thread.Upvotes == nil || thread.Downvotes == nil {
		zero := uint(0)
		thread.Upvotes, thread.Downvotes = &zero, &zero
	}

	err = h.validator.Struct(&thread)
	if err != nil {
		h.logger.Printf("Error in PostThread handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PostThread handler")
	}

	newThread, err := h.repository.CreateThread(thread)
	if err != nil {
		h.logger.Printf("Error in CreateThread repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateThread repository")
	}

	return c.JSON(http.StatusOK, &newThread)
}

// @Summary	Fetches a thread by its ID.
// @Tags		Threads
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Thread of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Thread}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/threads/{id} [get]
func (h *ForumHandler) GetThreadByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetThreadByID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetThreadByID handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	thread, err := h.repository.FindThreadByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Thread not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Thread not found!")
		}
		h.logger.Printf("Error in FindThreadByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindThreadByID repository")
	}

	return c.JSON(http.StatusOK, &thread)
}

// @Summary	Fetches all threads.
// @Tags		Threads
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Thread}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/threads [get]
func (h *ForumHandler) GetThreads(c echo.Context) error {
	threads, err := h.repository.FindThreads()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Thread not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Thread not found!")
		}
		h.logger.Printf("Error in FindFirstThread repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindThreads repository")
	}

	return c.JSON(http.StatusOK, &threads)
}

// @Summary	Updates an thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Thread of ID"
// @Param		Thread	body		models.Thread	true	"Update Thread"
// @Success	200		{object}	ResponseHTTP{data=models.Thread}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/threads/{id} [patch]
func (h *ForumHandler) PatchThread(c echo.Context) error {
	var thread models.Thread
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in PatchThread handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchThread handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	if err := c.Bind(&thread); err != nil {
		h.logger.Printf("Error in PatchThread handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchThread handler")
	}

	if thread.UpdatedAt == nil {
		currentTime := time.Now()
		thread.UpdatedAt = &currentTime
	}

	err = h.validator.Struct(&thread)
	if err != nil {
		h.logger.Printf("Error in PatchThread handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchThread handler")
	}

	updThread, err := h.repository.UpdateThread(int(id), thread)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Thread not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Thread not found!")
		}
		h.logger.Printf("Error in UpdateThread repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateThread repository")
	}

	return c.JSON(http.StatusOK, &updThread)
}

// @Summary	Deletes the specified thread.
// @Tags		Threads
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Thread of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/threads/{id} [delete]
func (h *ForumHandler) DeleteThread(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in DeleteThread handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in DeleteThread handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteThread(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Thread not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Thread not found!")
		}
		h.logger.Printf("Error in DeleteThread repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteThread repository")
	}

	return c.String(http.StatusOK, "Thread sucessfully deleted!")
}

// @Summary	Creates a new message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		Message	body		models.Message	true	"Create Message"
// @Success	200		{object}	ResponseHTTP{data=models.Message}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/messages [post]
func (h *ForumHandler) PostMessage(c echo.Context) error {
	var message models.Message

	err := c.Bind(&message)
	if err != nil {
		h.logger.Printf("Error in PostMessage handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostMessage handler")
	}
	message.Method = "POST"

	if message.CreatedAt == nil || message.UpdatedAt == nil {
		currentTime := time.Now()
		message.CreatedAt, message.UpdatedAt = &currentTime, &currentTime
	}
	if message.Upvotes == nil || message.Downvotes == nil {
		zero := uint(0)
		message.Upvotes, message.Downvotes = &zero, &zero
	}

	err = h.validator.Struct(&message)
	if err != nil {
		h.logger.Printf("Error in PostMessage handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PostMessage handler")
	}

	newMessage, err := h.repository.CreateMessage(message)
	if err != nil {
		h.logger.Printf("Error in CreateMessage repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateMessage repository")
	}

	return c.JSON(http.StatusOK, &newMessage)
}

// @Summary	Fetches a message by its ID.
// @Tags		Messages
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Message of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Message}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/messages/{id} [get]
func (h *ForumHandler) GetMessageByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetMessageByID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetMessageByID handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	message, err := h.repository.FindMessageByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Message not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Message not found!")
		}
		h.logger.Printf("Error in FindMessageByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindMessageByID repository")
	}

	return c.JSON(http.StatusOK, &message)
}

// @Summary	Fetches all messages.
// @Tags		Messages
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Message}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/messages [get]
func (h *ForumHandler) GetMessages(c echo.Context) error {
	messages, err := h.repository.FindMessages()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Message not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Message not found!")
		}
		h.logger.Printf("Error in FindFirstMessage repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindMessages repository")
	}

	return c.JSON(http.StatusOK, &messages)
}

// @Summary	Fetches all messages in the thread of ID.
// @Tags		Messages
// @Accept	text/plain
// @Produce	application/json
// @Param		threadID	path		int	true	"Get Messages of Thread ID"
// @Success	200	{object}	ResponseHTTP{data=[]models.Message}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/messages/thread/{threadID} [get]
func (h *ForumHandler) GetMessagesByThreadID(c echo.Context) error {
	p := c.Param("threadID")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetMessageByThreadID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetMessageByThreadID handler"+err.Error())
	}
	threadID, _ := strconv.ParseInt(p, 10, 64)

	messages, err := h.repository.FindMessagesByThreadID(int(threadID))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Message not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Message not found!")
		}
		h.logger.Printf("Error in FindFirstMessage repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindMessages repository")
	}

	return c.JSON(http.StatusOK, &messages)
}

// @Summary	Updates an message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Message of ID"
// @Param		Message	body		models.Message	true	"Update Message"
// @Success	200		{object}	ResponseHTTP{data=models.Message}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/messages/{id} [patch]
func (h *ForumHandler) PatchMessage(c echo.Context) error {
	var message models.Message
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in PatchMessage handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchMessage handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	if err := c.Bind(&message); err != nil {
		h.logger.Printf("Error in PatchMessage handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchMessage handler")
	}
	if message.UpdatedAt == nil {
		currentTime := time.Now()
		message.UpdatedAt = &currentTime
	}

	err = h.validator.Struct(&message)
	if err != nil {
		h.logger.Printf("Error in PatchMessage handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchMessage handler")
	}

	updMessage, err := h.repository.UpdateMessage(int(id), message)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Message not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Message not found!")
		}
		h.logger.Printf("Error in UpdateMessage repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateMessage repository")
	}

	return c.JSON(http.StatusOK, &updMessage)
}

// @Summary	Deletes the specified message.
// @Tags		Messages
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Message of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/messages/{id} [delete]
func (h *ForumHandler) DeleteMessage(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in DeleteMessage handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in DeleteMessage handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteMessage(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Message not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Message not found!")
		}
		h.logger.Printf("Error in DeleteMessage repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteMessage repository")
	}

	return c.String(http.StatusOK, "Message sucessfully deleted!")
}
