package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"strconv"

	_ "gamehangar/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AssetHandler struct {
	logger     echo.Logger
	repository AssetRepository
	validator  *validator.Validate
}

func NewAssetHandler(e *echo.Echo, repo AssetRepository, v *validator.Validate) *AssetHandler {
	return &AssetHandler{
		logger:     e.Logger,
		repository: repo,
		validator:  v,
	}
}

// @Summary	Creates a new asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		Asset	body		models.Asset	true	"Create Asset"
// @Success	201	{object}	models.Asset
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/assets [post]
func (h *AssetHandler) PostAsset(c echo.Context) error {
	var asset models.Asset

	err := c.Bind(&asset)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	asset.Method = "POST"

	err = h.validator.Struct(&asset)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	// TODO: Hook a Service to create links to the S3 bucket

	newAsset, err := h.repository.CreateAsset(asset)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateAsset repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newAsset)
}

// @Summary	Fetches a asset by its ID.
// @Tags		Assets
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Asset of ID"
// @Success	200	{object}	models.Asset
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/assets/{id} [get]
func (h *AssetHandler) GetAssetById(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetAssetByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	asset, err := h.repository.FindAssetByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindAssetByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &asset)
}

// @Summary	Fetches all assets.
// @Tags		Assets
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Success	200	{object}	models.Asset
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/assets [get]
func (h *AssetHandler) GetAssets(c echo.Context) error {
	var err error
	var assets *[]models.Asset

	tags := c.Request().URL.Query()["q"]
	if len(tags) != 0 {
		assets, err = h.repository.FindAssetsByQuery(&tags)
	} else {
		assets, err = h.repository.FindAssets()
	}
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindAssets repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &assets)
}

// @Summary	Updates an asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string			true	"Update Asset of ID"
// @Param		Asset	body		models.Asset	true	"Update Asset"
// @Success	200		{object}	models.Asset
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	409	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/assets/{id} [patch]
func (h *AssetHandler) PatchAsset(c echo.Context) error {
	var asset models.Asset
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = c.Bind(&asset)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	asset.Method = "PATCH"

	err = h.validator.Struct(&asset)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	updAsset, err := h.repository.UpdateAsset(int(id), asset)
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
				Message: "Error: unable to update the Asset due to an edit conflict, please try again!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusConflict, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateAsset repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updAsset)
}

// @Summary	Deletes the specified asset.
// @Tags		Assets
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Asset of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteAsset(int(id))
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
			Message: "Error in DeleteAsset repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Asset sucessfully deleted!")
}
