package v1

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"time"

	// _ "gamehangar/docs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AssetHandler struct {
	logger     echo.Logger
	repository AssetRepository
}

func NewAssetHandler(e *echo.Echo, repo AssetRepository) (*AssetHandler, error) {
	return &AssetHandler{
		logger:     e.Logger,
		repository: repo,
	}, nil
}

// @Summary	Creates a new asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		Asset	body		models.Asset	true	"Create Asset"
// @Success	200		{object}	ResponseHTTP{data=models.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/asset/ [post]
func (ah *AssetHandler) PostAsset(c echo.Context) error {
	var asset models.Asset

	err := c.Bind(&asset)
	if err != nil {
		ah.logger.Printf("Error in PostAsset handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PostAsset handler")
	}

	asset.ID = uuid.NewString()
	asset.CreatedAt = time.Now()
	// TODO: Hook a Service to create links to the S3 bucket

	newAsset, err := ah.repository.CreateAsset(asset)
	if err != nil {
		ah.logger.Printf("Error in CreateAsset operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateAsset repository")
	}

	return c.JSON(http.StatusOK, &newAsset)
}

// @Summary	Fetches a asset by its ID.
// @Tags		Assets
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Asset of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/asset/{id} [get]
func (ah *AssetHandler) GetAssetById(c echo.Context) error {
	id := c.Param("id")

	asset, err := ah.repository.FindAssetByID(id)
	if err.Error() == "Not Found" {
		ah.logger.Printf("Error: Asset not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Asset not found!")
	}
	if err != nil {
		ah.logger.Printf("Error in FindAssetByID operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in FindAssetByID operation")
	}

	return c.JSON(http.StatusOK, &asset)
}

// @Summary	Fetches all assets.
// @Tags		Assets
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/asset/ [get]
func (ah *AssetHandler) GetAssets(c echo.Context) error {
	assets, err := ah.repository.FindAssets()
	if err.Error() == "Not Found" {
		ah.logger.Printf("Error: Asset not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Asset not found!")
	}
	if err != nil {
		ah.logger.Printf("Error in FindFirstAsset operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in FindAssets operation")
	}

	return c.JSON(http.StatusOK, &assets)
}

// @Summary	Updates an asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		id	path		string	true	"Update Asset of ID"
// @Param		Asset	body		models.Asset	true	"Update Asset"
// @Success	200		{object}	ResponseHTTP{data=models.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/asset/{id} [patch]
func (ah *AssetHandler) PatchAsset(c echo.Context) error {
	var asset models.Asset
	id := c.Param("id")

	err := c.Bind(&asset)
	if err != nil {
		ah.logger.Printf("Error in PatchAsset handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PatchAsset handler")
	}

	updAsset, err := ah.repository.UpdateAsset(id, asset)
	if err.Error() == "Not Found" {
		ah.logger.Printf("Error: Asset not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Asset not found!")
	}
	if err != nil {
		ah.logger.Printf("Error in UpdateAsset repository \n%s", err)
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
// @Router		/v1/asset/{id} [delete]
func (ah *AssetHandler) DeleteAsset(c echo.Context) error {
	id := c.Param("id")

	err := ah.repository.DeleteAsset(id)
	if err.Error() == "Not Found" {
		ah.logger.Printf("Error: Asset not found!\n%s", err)
		return c.String(http.StatusNotFound, "Error: Asset not found!")
	}
	if err != nil {
		ah.logger.Printf("Error in DeleteAsset repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteAsset repository")
	}

	return c.String(http.StatusOK, "Asset sucessfully deleted!")
}
