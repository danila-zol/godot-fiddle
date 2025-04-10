package v1

import (
	"encoding/json"
	"gamehangar/internal/domain/models"
	"log"
	"net/http"
	"time"

	// _ "gamehangar/docs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type AssetHandler struct {
	echo       *echo.Echo
	repository AssetRepository
}

func NewAssetHandler(e *echo.Echo, repo AssetRepository) (*AssetHandler, error) {
	return &AssetHandler{
		echo:       e,
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
// @Router		/assets/ [post]
func (ah *AssetHandler) postAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in PostAsset handler\n"))
		log.Printf("Error in PostAsset handler \n%s", err)
		return
	}

	asset.ID = uuid.NewString()
	asset.CreatedAt = time.Now()
	// TODO: Hook a Service to create links to the S3 bucket

	newAsset, err := ah.repository.CreateAsset(asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in CreateAsset operation\n"))
		log.Printf("Error in CreateAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newAsset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in postAsset operation\n"))
		log.Printf("Error in postAsset operation \n%s", err)
		return
	}
}

// @Summary	Fetches a asset by its ID.
// @Tags		Assets
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Asset of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/assets/{id} [get]
func (ah *AssetHandler) getAssetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	asset, err := ah.repository.FindAssetByID(id)
	if err.Error() == "Not Found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error: Asset not found!\n"))
		log.Printf("Error: Asset not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in FindFirstAsset operation\n"))
		log.Printf("Error in FindFirstAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in getAssetById operation\n"))
		log.Printf("Error in getAssetById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all assets.
// @Tags		Assets
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/assets/ [get]
func (ah *AssetHandler) getAssets(w http.ResponseWriter, r *http.Request) {
	asset, err := ah.repository.FindAssets()
	if err.Error() == "Not Found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error: Assets not found!\n"))
		log.Printf("Error: Assets not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in FindAssets operation\n"))
		log.Printf("Error in FindAssets operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in getAssets operation\n"))
		log.Printf("Error in getAssets operation \n%s", err)
		return
	}
}

// @Summary	Updates an asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		id	path		string	true	"Delete Asset of ID"
// @Param		Asset	body		models.Asset	true	"Update Asset"
// @Success	200		{object}	ResponseHTTP{data=models.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/assets/ [patch]
func (ah *AssetHandler) patchAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	id := r.PathValue("id")
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in patchAsset handler\n"))
		log.Printf("Error in patchAsset handler \n%s", err)
		return
	}

	updAsset, err := ah.repository.UpdateAsset(id, asset)
	if err.Error() == "Not Found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error: Assets not found!\n"))
		log.Printf("Error: Assets not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in UpdateAsset operation\n"))
		log.Printf("Error in UpdateAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updAsset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in patchAsset handler\n"))
		log.Printf("Error in patchAsset handler \n%s", err)
		return
	}
}

// @Summary	Deletes the specified asset.
// @Tags		Assets
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Asset of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/assets/{id} [delete]
func (ah *AssetHandler) deleteAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := ah.repository.DeleteAsset(id)
	if err.Error() == "Not Found" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error: Asset not found!\n"))
		log.Printf("Error: Asset not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in FindFirstAsset operation\n"))
		log.Printf("Error in FindFirstAsset operation \n%s", err)
		return
	}
}
