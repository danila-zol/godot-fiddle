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
// @Success	201	{object}	models.Demo
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/demos [post]
func (h *DemoHandler) PostDemo(c echo.Context) error {
	var demo models.Demo

	err := c.Bind(&demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	demo.Method = "POST"

	err = h.validator.Struct(&demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	demo.ThreadID, err = h.syncer.PostThread(demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	newDemo, err := h.repository.CreateDemo(demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateDemo repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newDemo)
}

// @Summary	Fetches a demo by its ID.
// @Tags		Demos
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Demo of ID"
// @Success	200	{object}	models.Demo
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/demos/{id} [get]
func (h *DemoHandler) GetDemoById(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetDemoByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	demo, err := h.repository.FindDemoByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindDemoByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &demo)
}

// @Summary	Fetches all demos.
// @Tags		Demos
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Param		l	query		int	false	"Record number limit"
// @Param		o	query		string	false	"Record ordering. Default newest updated" Enums(newest-updated, highest-rated, most-views)
// @Success	200	{object}	models.Demo
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/demos [get]
func (h *DemoHandler) GetDemos(c echo.Context) error {
	var (
		err   error
		limit uint64
		order string
		demos *[]models.Demo
	)

	l := c.Request().URL.Query()["l"]
	if l != nil {
		err = h.validator.Var(l[0], "omitnil,number,min=0")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetDemos repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		limit, err = strconv.ParseUint(l[0], 10, 64)
	}
	tags := c.Request().URL.Query()["q"]

	o := c.Request().URL.Query()["o"]
	if o != nil {
		err = h.validator.Var(o[0], `oneof=newest-updated highest-rated most-views`)
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
		order = "newest-updated"
	}

	demos, err = h.repository.FindDemos(tags, limit, order)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindDemos repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &demos)
}

// @Summary	Updates a demo.
// @Tags		Demos
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		int		true	"Update Demo of ID"
// @Param		Demo	body		models.Demo	true	"Update Demo"
// @Success	200		{object}	models.Demo
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/demos/{id} [patch]
func (h *DemoHandler) PatchDemo(c echo.Context) error {
	var demo models.Demo

	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	demo.Method = "PATCH"

	err = h.validator.Struct(&demo)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	err = h.syncer.PatchThread(int(id), demo)
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
			Message: "Error in UpdateDemo repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	updDemo, err := h.repository.UpdateDemo(int(id), demo)
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
			Message: "Error in UpdateDemo repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updDemo)
}

// @Summary	Deletes the specified demo.
// @Tags		Demos
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		int	true	"Delete Demo of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/demos/{id} [delete]
func (h *DemoHandler) DeleteDemo(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteDemo handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteDemo(int(id))
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
			Message: "Error in DeleteDemo repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Demo successfully deleted!")
}
