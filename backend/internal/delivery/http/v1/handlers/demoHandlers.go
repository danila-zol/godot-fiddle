package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"time"

	_ "gamehangar/docs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DemoHandler struct {
	logger     echo.Logger
	repository DemoRepository
	syncer     ThreadSyncer
}

type ThreadSyncer interface {
	PostThread(demo models.Demo) (*string, error)
	PatchThread(demo models.Demo) error
}

func NewDemoHandler(e *echo.Echo, repo DemoRepository, s ThreadSyncer) *DemoHandler {
	return &DemoHandler{
		logger:     e.Logger,
		repository: repo,
		syncer:     s,
	}
}

// @Summary	Creates a new demo.
// @Tags		Demos
// @Accept		application/json
// @Produce	application/json
// @Param		Demo	body		models.Demo	true	"Create Demo"
// @Success	200		{object}	ResponseHTTP{data=models.Demo}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/demos [post]
func (h *DemoHandler) PostDemo(c echo.Context) error {
	var demo models.Demo

	err := c.Bind(&demo)
	if err != nil {
		h.logger.Printf("Error in PostDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostDemo handler")
	}

	if demo.ID == nil {
		demoID := uuid.NewString()
		demo.ID = &demoID
	}
	if demo.CreatedAt == nil || demo.UpdatedAt == nil {
		currentTime := time.Now()
		demo.CreatedAt, demo.UpdatedAt = &currentTime, &currentTime
	}
	if demo.Upvotes == nil || demo.Downvotes == nil {
		zero := uint(0)
		demo.Upvotes, demo.Downvotes = &zero, &zero
	}
	if demo.Tags == nil {
		empty := make([]string, 0)
		demo.Tags = &empty
	}

	demo.ThreadID, err = h.syncer.PostThread(demo)
	if err != nil {
		h.logger.Printf("Error in PostDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostDemo handler")
	}

	newDemo, err := h.repository.CreateDemo(demo)
	if err != nil {
		h.logger.Printf("Error in CreateDemo repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateDemo repository")
	}

	return c.JSON(http.StatusOK, &newDemo)
}

// @Summary	Fetches a demo by its ID.
// @Tags		Demos
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Demo of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Demo}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [get]
func (h *DemoHandler) GetDemoById(c echo.Context) error {
	id := c.Param("id")

	demo, err := h.repository.FindDemoByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in FindDemoByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateDemo repository")
	}

	return c.JSON(http.StatusOK, &demo)
}

// @Summary	Fetches all demos.
// @Tags		Demos
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Demo}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos [get]
func (h *DemoHandler) GetDemos(c echo.Context) error {
	demos, err := h.repository.FindDemos()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in FindDemos operation: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindDemos operation")
	}

	return c.JSON(http.StatusOK, &demos)
}

// @Summary	Updates a demo.
// @Tags		Demos
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string		true	"Update Demo of ID"
// @Param		Demo	body		models.Demo	true	"Update Demo"
// @Success	200		{object}	ResponseHTTP{data=models.Demo}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [patch]
func (h *DemoHandler) PatchDemo(c echo.Context) error {
	var demo models.Demo

	id := c.Param("id")

	err := c.Bind(&demo)
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchDemo handler")
	}

	if demo.UpdatedAt == nil {
		currentTime := time.Now()
		demo.UpdatedAt = &currentTime
	}

	err = h.syncer.PatchThread(demo)
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchDemo handler")
	}

	updDemo, err := h.repository.UpdateDemo(id, demo)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in UpdateDemo repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateDemo repository")
	}

	return c.JSON(http.StatusOK, &updDemo)
}

// @Summary	Deletes the specified demo.
// @Tags		Demos
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Demo of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [delete]
func (h *DemoHandler) DeleteDemo(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteDemo(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in DeleteDemo repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteDemo repository")
	}

	return c.String(http.StatusOK, "Demo successfully deleted!")
}
