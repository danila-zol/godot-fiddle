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
}

func NewDemoHandler(e *echo.Echo, repo DemoRepository) (*DemoHandler, error) {
	return &DemoHandler{
		logger:     e.Logger,
		repository: repo,
	}, nil
}

// @Summary	Creates a new demo.
// @Tags		Demos
// @Accept		application/json
// @Produce	application/json
// @Param		Demo	body		models.Demo	true	"Create Demo"
// @Success	200		{object}	ResponseHTTP{data=models.Demo}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/demo/ [post]
func (h *DemoHandler) PostDemo(c echo.Context) error {
	var demo models.Demo
	err := c.Bind(&demo)
	if err != nil {
		h.logger.Printf("Error in PostAsset handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PostAsset handler")
	}

	demo.ID = uuid.NewString()
	demo.CreatedAt, demo.UpdatedAt = time.Now(), time.Now()
	demo.Upvotes, demo.Downvotes = 0, 0

	// TODO: Service that sends a POST request to create a new Thread and a new Message

	newDemo, err := h.repository.CreateDemo(demo)
	if err != nil {
		h.logger.Printf("Error in CreateAsset operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateAsset repository")
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
// @Router		/v1/demo/{id} [get]
func (h *DemoHandler) GetDemoById(c echo.Context) error {
	id := c.Param("id")

	demo, err := h.repository.FindDemoByID(id)
	if err.Error() == "Not Found" {
		h.logger.Printf("Error: Demos not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Demos not found!")
	}
	if err != nil {
		h.logger.Printf("Error in FindDemoByID repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateAsset repository")
	}

	return c.JSON(http.StatusOK, &demo)
}

// @Summary	Fetches all demos.
// @Tags		Demos
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Demo}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/demo/ [get]
func (h *DemoHandler) GetDemos(c echo.Context) error {
	demos, err := h.repository.FindDemos()
	if err.Error() == "Not Found" {
		h.logger.Printf("Error: Demos not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Demos not found!")
	}
	if err != nil {
		h.logger.Printf("Error in FindDemos operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in FindDemos operation")
	}

	return c.JSON(http.StatusOK, &demos)
}

// @Summary	Updates a demo.
// @Tags		Demos
// @Accept		application/json
// @Produce	application/json
// @Param		id	path		string	true	"Update Demo of ID"
// @Param		Demo	body		models.Demo	true	"Update Demo"
// @Success	200		{object}	ResponseHTTP{data=models.Demo}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/demo/{id} [patch]
func (h *DemoHandler) PatchDemo(c echo.Context) error {
	var demo models.Demo

	id := c.Param("id")

	err := c.Bind(&demo)
	if err != nil {
		h.logger.Printf("Error in PatchDemo handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PatchDemo handler")
	}
	demo.UpdatedAt = time.Now()

	// TODO: Update Thread and its initial Message when its Demo updates

	updDemo, err := h.repository.UpdateDemo(id, demo)
	if err.Error() == "Not Found" {
		h.logger.Printf("Error: Demo not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Demo not found!")
	}
	if err != nil {
		h.logger.Printf("Error in UpdateDemo repository \n%s", err)
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
// @Router		/v1/demo/{id} [delete]
func (h *DemoHandler) DeleteDemo(c echo.Context) error {
	id := c.Param("id")

	// TODO: Delete the respective Thread and the initial Message?

	err := h.repository.DeleteDemo(id)
	if err.Error() == "Not Found" {
		h.logger.Printf("Error: Demo not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Demo not found!")
	}
	if err != nil {
		h.logger.Printf("Error in DeleteDemo repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteDemo repository")
	}

	return c.String(http.StatusOK, "Demo successfully deleted!")
}
