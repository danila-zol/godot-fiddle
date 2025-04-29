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
// @Success	200		{object}	ResponseHTTP{data=models.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/assets [post]
func (h *AssetHandler) PostAsset(c echo.Context) error {
	var asset models.Asset

	err := c.Bind(&asset)
	if err != nil {
		h.logger.Printf("Error in PostAsset handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostAsset handler")
	}
	asset.Method = "POST"

	err = h.validator.Struct(&asset)
	if err != nil {
		h.logger.Printf("Error in PostAsset handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PostAsset handler")
	}

	// TODO: Hook a Service to create links to the S3 bucket

	newAsset, err := h.repository.CreateAsset(asset)
	if err != nil {
		h.logger.Printf("Error in CreateAsset repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateAsset repository")
	}

	return c.JSON(http.StatusCreated, &newAsset)
}

// @Summary	Fetches a asset by its ID.
// @Tags		Assets
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		int	true	"Get Asset of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/assets/{id} [get]
func (h *AssetHandler) GetAssetById(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in GetAssetByID handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in GetAssetByID handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	asset, err := h.repository.FindAssetByID(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Asset not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in FindAssetByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindAssetByID repository")
	}

	return c.JSON(http.StatusOK, &asset)
}

// @Summary	Fetches all assets.
// @Tags		Assets
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Success	200	{object}	ResponseHTTP{data=[]models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
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
			h.logger.Printf("Error: Asset not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in FindFirstAsset repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindAssets repository")
	}

	return c.JSON(http.StatusOK, &assets)
}

// @Summary	Updates an asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string			true	"Update Asset of ID"
// @Param		Asset	body		models.Asset	true	"Update Asset"
// @Success	200		{object}	ResponseHTTP{data=models.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/assets/{id} [patch]
func (h *AssetHandler) PatchAsset(c echo.Context) error {
	var asset models.Asset
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in PatchAsset handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchAsset handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	if err = c.Bind(&asset); err != nil {
		h.logger.Printf("Error in PatchAsset handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchAsset handler")
	}
	asset.Method = "PATCH"

	err = h.validator.Struct(&asset)
	if err != nil {
		h.logger.Printf("Error in PatchAsset handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in PatchAsset handler")
	}

	updAsset, err := h.repository.UpdateAsset(int(id), asset)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Asset not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		} else if err == h.repository.ConflictErr() {
			h.logger.Printf("Error: unable to update the Asset due to an edit conflict, please try again!")
			return c.String(http.StatusConflict, "Error: unable to update the Asset due to an edit conflict, please try again!")
		}
		h.logger.Printf("Error in UpdateAsset repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateAsset repository")
	}

	return c.JSON(http.StatusOK, &updAsset)
}

// @Summary	Deletes the specified asset.
// @Tags		Assets
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Asset of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c echo.Context) error {
	p := c.Param("id")
	err := h.validator.Var(p, "required,number")
	if err != nil {
		h.logger.Printf("Error in DeleteAsset handler: %s", err)
		return c.String(http.StatusUnprocessableEntity, "Error in DeleteAsset handler"+err.Error())
	}
	id, _ := strconv.ParseInt(p, 10, 64)

	err = h.repository.DeleteAsset(int(id))
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Asset not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in DeleteAsset repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteAsset repository")
	}

	return c.String(http.StatusOK, "Asset sucessfully deleted!")
}
