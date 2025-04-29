package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"strconv"

	_ "gamehangar/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DemoHandler struct {
	logger     echo.Logger
	repository DemoRepository
	validator  *validator.Validate
	syncer     ThreadSyncer
}

type ThreadSyncer interface {
	PostThread(demo models.Demo) (*int, error)
	PatchThread(demoID int, demo models.Demo) error
}

func NewDemoHandler(e *echo.Echo, repo DemoRepository, v *validator.Validate, s ThreadSyncer) *DemoHandler {
	return &DemoHandler{
		logger:     e.Logger,
		repository: repo,
		validator:  v,
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
	demo.Method = "POST"

	err = h.validator.Struct(&demo)
	if err != nil {
		h.logger.Printf("Error in PostDemo handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PostDemo handler")
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
// @Param		id	path		int	true	"Get Demo of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Demo}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [get]
func (h *DemoHandler) GetDemoById(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetDemoByID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetDemoByID handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	demo, err := h.repository.FindDemoByID(int(id))
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
// @Param		q	query		[]string	false	"Keyword Query"
// @Success	200	{object}	ResponseHTTP{data=[]models.Demo}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos [get]
func (h *DemoHandler) GetDemos(c echo.Context) error {
	var err error
	var demos *[]models.Demo

	tags := c.Request().URL.Query()["q"]
	if len(tags) != 0 {
		demos, err = h.repository.FindDemosByQuery(&tags)
	} else {
		demos, err = h.repository.FindDemos()
	}
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Demos not found!")
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
// @Param		id		path		int		true	"Update Demo of ID"
// @Param		Demo	body		models.Demo	true	"Update Demo"
// @Success	200		{object}	ResponseHTTP{data=models.Demo}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [patch]
func (h *DemoHandler) PatchDemo(c echo.Context) error {
	var demo models.Demo

	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchDemo handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&demo)
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchDemo handler")
	}
	demo.Method = "PATCH"

	err = h.validator.Struct(&demo)
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchDemo handler")
	}

	err = h.syncer.PatchThread(int(id), demo)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Demos not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in PatchDemo handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchDemo handler")
	}

	updDemo, err := h.repository.UpdateDemo(int(id), demo)
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
// @Param		id	path		int	true	"Delete Demo of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demos/{id} [delete]
func (h *DemoHandler) DeleteDemo(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in DeleteDemo handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in DeleteDemo handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteDemo(int(id))
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
