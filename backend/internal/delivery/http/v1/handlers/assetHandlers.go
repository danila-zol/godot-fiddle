package handlers

import (
	"gamehangar/internal/domain/models"
	"mime/multipart"
	"net/http"
	"strconv"

	_ "gamehangar/docs"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AssetHandler struct {
	logger         echo.Logger
	repository     AssetRepository
	validator      *validator.Validate
	objectUploader ObjectUploader
}

func NewAssetHandler(e *echo.Echo, repo AssetRepository, v *validator.Validate, o ObjectUploader) *AssetHandler {
	return &AssetHandler{
		logger:         e.Logger,
		repository:     repo,
		validator:      v,
		objectUploader: o,
	}
}

//	@Summary	Creates a new asset.
//	@Tags		Assets
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		Asset			formData	models.Asset	true	"Create Asset"
//	@param		assetFile		formData	file			true	"Asset project file"
//	@param		assetThumbnail	formData	file			true	"Asset thumbnail"
//	@Success	201				{object}	models.Asset
//	@Failure	400				{object}	HTTPError
//	@Failure	403				{object}	HTTPError
//	@Failure	404				{object}	HTTPError
//	@Failure	413				{object}	HTTPError
//	@Failure	422				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Router		/v1/assets [post]
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

	assetFormFile, err := c.FormFile("assetFile")
	if assetFormFile == nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Missing asset project file: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	assetMultipartFile, err := assetFormFile.Open()
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error uploading file! Please try again",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	defer assetMultipartFile.Close()

	err = h.objectUploader.CheckFileSize(assetFormFile.Size, c.Get("userTier").(string))
	if err != nil {
		if err == h.objectUploader.ObjectTooLargeErr() {
			e := HTTPError{
				Code:    http.StatusRequestEntityTooLarge,
				Message: "Error in PostAsset handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusRequestEntityTooLarge, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in PostAsset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}
	thumbnailFormFile, err := c.FormFile("assetThumbnail")
	if thumbnailFormFile == nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Missing asset thumnail file: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	thumbnailMultipartFile, err := thumbnailFormFile.Open()
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error uploading file! Please try again",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	defer thumbnailMultipartFile.Close()
	err = h.objectUploader.CheckFileSize(thumbnailFormFile.Size, "picture")
	if err != nil {
		if err == h.objectUploader.ObjectTooLargeErr() {
			e := HTTPError{
				Code:    http.StatusRequestEntityTooLarge,
				Message: "Error in Postasset handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusRequestEntityTooLarge, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in Postasset handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	newAsset, err := h.repository.CreateAsset(asset, assetMultipartFile, thumbnailMultipartFile)
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

//	@Summary	Fetches a asset by its ID.
//	@Tags		Assets
//	@Accept		text/plain
//	@Produce	application/json
//	@Param		id	path		int	true	"Get Asset of ID"
//	@Success	200	{object}	models.Asset
//	@Failure	400	{object}	HTTPError
//	@Failure	404	{object}	HTTPError
//	@Failure	422	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/v1/assets/{id} [get]
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

//	@Summary	Fetches all assets.
//	@Tags		Assets
//	@Produce	application/json
//	@Param		q	query		[]string	false	"Keyword Query"
//	@Param		l	query		int			false	"Record number limit"
//	@Param		o	query		string		false	"Record ordering. Default newest updated"	Enums(newest-updated, highest-rated, most-views)
//	@Success	200	{object}	models.Asset
//	@Failure	400	{object}	HTTPError
//	@Failure	404	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/v1/assets [get]
func (h *AssetHandler) GetAssets(c echo.Context) error {
	var (
		err    error
		limit  uint64
		order  string
		assets *[]models.Asset
	)

	l := c.Request().URL.Query()["l"]
	if l != nil {
		err = h.validator.Var(l[0], "omitnil,number,min=0")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetAssets repository: " + err.Error(),
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

	assets, err = h.repository.FindAssets(tags, limit, order)
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

//	@Summary	Updates an asset.
//	@Tags		Assets
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		id				path		string			true	"Update Asset of ID"
//	@Param		Asset			formData	models.Asset	true	"Update Asset"
//	@param		assetFile		formData	file			false	"Asset project file"
//	@param		assetThumbnail	formData	file			false	"Asset thumbnail"
//	@Success	200				{object}	models.Asset
//	@Failure	400				{object}	HTTPError
//	@Failure	403				{object}	HTTPError
//	@Failure	404				{object}	HTTPError
//	@Failure	409				{object}	HTTPError
//	@Failure	413				{object}	HTTPError
//	@Failure	422				{object}	HTTPError
//	@Failure	500				{object}	HTTPError
//	@Router		/v1/assets/{id} [patch]
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

	var assetMultipartFile, thumbnailMultipartFile multipart.File
	assetFormFile, err := c.FormFile("assetFile")
	if assetFormFile != nil {
		assetMultipartFile, err = assetFormFile.Open()
		if err != nil {
			e := HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Error uploading file! Please try again",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusBadRequest, &e)
		}
		defer assetMultipartFile.Close()

		err = h.objectUploader.CheckFileSize(assetFormFile.Size, c.Get("userTier").(string))
		if err != nil {
			if err == h.objectUploader.ObjectTooLargeErr() {
				e := HTTPError{
					Code:    http.StatusRequestEntityTooLarge,
					Message: "Error in PatchAsset handler: " + err.Error(),
				}
				h.logger.Print(&e)
				return c.JSON(http.StatusRequestEntityTooLarge, &e)
			}
			e := HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Error in PatchAsset handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusInternalServerError, &e)
		}
	}

	thumbnailFormFile, err := c.FormFile("assetThumbnail")
	if thumbnailFormFile != nil {
		thumbnailMultipartFile, err = thumbnailFormFile.Open()
		if err != nil {
			e := HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Error uploading file! Please try again",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusBadRequest, &e)
		}
		defer thumbnailMultipartFile.Close()

		err = h.objectUploader.CheckFileSize(thumbnailFormFile.Size, "picture")
		if err != nil {
			if err == h.objectUploader.ObjectTooLargeErr() {
				e := HTTPError{
					Code:    http.StatusRequestEntityTooLarge,
					Message: "Error in PatchAsset handler: " + err.Error(),
				}
				h.logger.Print(&e)
				return c.JSON(http.StatusRequestEntityTooLarge, &e)
			}
			e := HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Error in PatchAsset handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusInternalServerError, &e)
		}
	}

	updAsset, err := h.repository.UpdateAsset(int(id), asset, assetMultipartFile, thumbnailMultipartFile)
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
			Message: "Error in UpdateAsset repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updAsset)
}

//	@Summary	Deletes the specified asset.
//	@Tags		Assets
//	@Accept		text/plain
//	@Produce	text/plain
//	@Param		id	path		string	true	"Delete Asset of ID"
//	@Success	200	{string}	string
//	@Failure	403	{object}	HTTPError
//	@Failure	404	{object}	HTTPError
//	@Failure	422	{object}	HTTPError
//	@Failure	500	{object}	HTTPError
//	@Router		/v1/assets/{id} [delete]
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
