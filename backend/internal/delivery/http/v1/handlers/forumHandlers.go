package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"strconv"

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
// @Success	201	{object}	models.Topic
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/topics [post]
func (h *ForumHandler) PostTopic(c echo.Context) error {
	var topic models.Topic

	err := c.Bind(&topic)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}

	err = h.validator.Struct(&topic)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	newTopic, err := h.repository.CreateTopic(topic)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateTopic repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newTopic)
}

// @Summary	Fetches a topic by its ID.
// @Tags		Topics
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Topic of ID"
// @Success	200	{object}	models.Topic
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/topics/{id} [get]
func (h *ForumHandler) GetTopicByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetTopicByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	topic, err := h.repository.FindTopicByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindTopicByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &topic)
}

// @Summary	Fetches all topics.
// @Tags		Topics
// @Produce	application/json
// @Success	200	{object}	models.Topic
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/topics [get]
func (h *ForumHandler) GetTopics(c echo.Context) error {
	topics, err := h.repository.FindTopics()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindTopics repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &topics)
}

// @Summary	Updates an topic.
// @Tags		Topics
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Topic of ID"
// @Param		Topic	body		models.Topic	true	"Update Topic"
// @Success	200		{object}	models.Topic
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	409	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/topics/{id} [patch]
func (h *ForumHandler) PatchTopic(c echo.Context) error {
	var topic models.Topic
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&topic)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	topic.Method = "PATCH"

	err = h.validator.Struct(&topic)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	updTopic, err := h.repository.UpdateTopic(int(id), topic)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		} else if err == h.repository.ConflictErr() {
			e := HTTPError{
				Code:    http.StatusConflict,
				Message: "Error: unable to update the record due to an edit conflict, please try again!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusConflict, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateTopic repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updTopic)
}

// @Summary	Deletes the specified topic.
// @Tags		Topics
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Topic of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/topics/{id} [delete]
func (h *ForumHandler) DeleteTopic(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteTopic handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteTopic(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteTopic repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Topic sucessfully deleted!")
}

// @Summary	Creates a new thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		Thread	body		models.Thread	true	"Create Thread"
// @Success	201	{object}	models.Thread
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/threads [post]
func (h *ForumHandler) PostThread(c echo.Context) error {
	var thread models.Thread

	err := c.Bind(&thread)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	thread.Method = "POST"

	err = h.validator.Struct(&thread)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	newThread, err := h.repository.CreateThread(thread)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateThread repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newThread)
}

// @Summary	Fetches a thread by its ID.
// @Tags		Threads
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Thread of ID"
// @Success	200	{object}	models.Thread
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/threads/{id} [get]
func (h *ForumHandler) GetThreadByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetThreadByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	thread, err := h.repository.FindThreadByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindThreadByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &thread)
}

// @Summary	Fetches all threads.
// @Tags		Threads
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Param		l	query		int	false	"Record number limit"
// @Param		o	query		string	false	"Record ordering. Default newest updated" Enums(newestUpdated, highestRated, mostViews)
// @Success	200	{object}	models.Thread
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/threads [get]
func (h *ForumHandler) GetThreads(c echo.Context) error {
	var (
		err     error
		limit   uint64
		order   string
		threads *[]models.Thread
	)

	l := c.Request().URL.Query()["l"]
	if l != nil {
		err = h.validator.Var(l[0], "omitnil,number,min=0")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetThreads repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		limit, err = strconv.ParseUint(l[0], 10, 64)
	}
	tags := c.Request().URL.Query()["q"]

	o := c.Request().URL.Query()["o"]
	if o != nil {
		err = h.validator.Var(o[0], `oneof=newestUpdated highestRated mostViews`)
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetAssets repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		order = o[0]
	} else {
		order = "newestUpdated"
	}

	threads, err = h.repository.FindThreads(tags, limit, order)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindThreads repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &threads)
}

// @Summary	Updates an thread.
// @Tags		Threads
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Thread of ID"
// @Param		Thread	body		models.Thread	true	"Update Thread"
// @Success	200		{object}	models.Thread
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/threads/{id} [patch]
func (h *ForumHandler) PatchThread(c echo.Context) error {
	var thread models.Thread
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&thread)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}

	err = h.validator.Struct(&thread)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	updThread, err := h.repository.UpdateThread(int(id), thread)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateThread repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updThread)
}

// @Summary	Deletes the specified thread.
// @Tags		Threads
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Thread of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/threads/{id} [delete]
func (h *ForumHandler) DeleteThread(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteThread handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteThread(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteThread repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Thread sucessfully deleted!")
}

// @Summary	Creates a new message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		Message	body		models.Message	true	"Create Message"
// @Success	201	{object}	models.Message
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages [post]
func (h *ForumHandler) PostMessage(c echo.Context) error {
	var message models.Message

	err := c.Bind(&message)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	message.Method = "POST"

	err = h.validator.Struct(&message)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	newMessage, err := h.repository.CreateMessage(message)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateMessage repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newMessage)
}

// @Summary	Fetches a message by its ID.
// @Tags		Messages
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Message of ID"
// @Success	200	{object}	models.Message
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages/{id} [get]
func (h *ForumHandler) GetMessageByID(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetMessageByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	message, err := h.repository.FindMessageByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindMessageByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &message)
}

// @Summary	Fetches all messages.
// @Tags		Messages
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Param		l	query		int	false	"Record number limit"
// @Param		o	query		string	false	"Record ordering. Default newest updated" Enums(newestUpdated, highestRated, mostViews)
// @Success	200	{object}	models.Message
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages [get]
func (h *ForumHandler) GetMessages(c echo.Context) error {
	var (
		err      error
		limit    uint64
		order    string
		messages *[]models.Message
	)

	l := c.Request().URL.Query()["l"]
	if l != nil {
		err = h.validator.Var(l[0], "omitnil,number,min=0")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetMessages repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		limit, err = strconv.ParseUint(l[0], 10, 64)
	}
	tags := c.Request().URL.Query()["q"]

	o := c.Request().URL.Query()["o"]
	if o != nil {
		err = h.validator.Var(o[0], `oneof=newestUpdated highestRated mostViews`)
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetAssets repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		order = o[0]
	} else {
		order = "newestUpdated"
	}

	messages, err = h.repository.FindMessages(tags, limit, order)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindMessages repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &messages)
}

// @Summary	Fetches all messages in the thread of ID.
// @Tags		Messages
// @Accept	text/plain
// @Produce	application/json
// @Param		threadID	path		int	true	"Get Messages of Thread ID"
// @Success	200	{object}	models.Message
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages/thread/{threadID} [get]
func (h *ForumHandler) GetMessagesByThreadID(c echo.Context) error {
	p := c.Param("threadID")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetMessageByMessageID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	threadID, _ := strconv.ParseInt(p, 10, 64)

	messages, err := h.repository.FindMessagesByThreadID(int(threadID))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindMessagesByThreadID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &messages)
}

// @Summary	Updates an message.
// @Tags		Messages
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int			true	"Update Message of ID"
// @Param		Message	body		models.Message	true	"Update Message"
// @Success	200		{object}	models.Message
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages/{id} [patch]
func (h *ForumHandler) PatchMessage(c echo.Context) error {
	var message models.Message
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&message)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}

	err = h.validator.Struct(&message)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	updMessage, err := h.repository.UpdateMessage(int(id), message)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateMessage repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updMessage)
}

// @Summary	Deletes the specified message.
// @Tags		Messages
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Message of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/messages/{id} [delete]
func (h *ForumHandler) DeleteMessage(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteMessage handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteMessage(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteMessage repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Message sucessfully deleted!")
}
