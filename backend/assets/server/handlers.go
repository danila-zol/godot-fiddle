package server

import (
	"encoding/json"
	operations "game-hangar/database"
	"log"
	"net/http"
	"time"

	_ "game-hangar/docs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// @Summary	Creates a new asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		Asset	body		database.Asset	true	"Create Asset"
// @Success	200		{object}	ResponseHTTP{data=database.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/assets/ [post]
func postAsset(w http.ResponseWriter, r *http.Request) {
	var asset operations.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in PostAsset handler\n"))
		log.Printf("Error in PostAsset handler \n%s", err)
		return
	}

	asset.ID = "asset_" + uuid.NewString()
	asset.Created_at = time.Now()

	newAsset, err := operations.CreateAsset(asset)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in CreateAsset operation\n"))
		log.Printf("Error in CreateAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(newAsset)
	if err != nil {
		w.WriteHeader(500)
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
// @Success	200	{object}	ResponseHTTP{data=database.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/assets/{id} [get]
func getAssetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	asset, err := operations.FindFirstAsset(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Asset not found!\n"))
		log.Printf("Error: Asset not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstAsset operation\n"))
		log.Printf("Error in FindFirstAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(asset)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getAssetById operation\n"))
		log.Printf("Error in getAssetById operation \n%s", err)
		return
	}
}

// @Summary	Fetches all assets.
// @Tags		Assets
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]database.Asset}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/assets/ [get]
func getAssets(w http.ResponseWriter, r *http.Request) {
	asset, err := operations.FindAssets()
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Assets not found!\n"))
		log.Printf("Error: Assets not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindAssets operation\n"))
		log.Printf("Error in FindAssets operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(asset)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in getAssets operation\n"))
		log.Printf("Error in getAssets operation \n%s", err)
		return
	}
}

// @Summary	Updates a asset.
// @Tags		Assets
// @Accept		application/json
// @Produce	application/json
// @Param		Asset	body		database.Asset	true	"Update Asset"
// @Success	200		{object}	ResponseHTTP{data=database.Asset}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/assets/ [patch]
func patchAsset(w http.ResponseWriter, r *http.Request) {
	var asset operations.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error in patchAsset handler\n"))
		log.Printf("Error in patchAsset handler \n%s", err)
		return
	}

	updAsset, err := operations.UpdateAsset(asset)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Assets not found!\n"))
		log.Printf("Error: Assets not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in UpdateAsset operation\n"))
		log.Printf("Error in UpdateAsset operation \n%s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(updAsset)
	if err != nil {
		w.WriteHeader(500)
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
func deleteAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := operations.DeleteAsset(id)
	if err == pgx.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("Error: Asset not found!\n"))
		log.Printf("Error: Asset not found!\n%s", err)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error in FindFirstAsset operation\n"))
		log.Printf("Error in FindFirstAsset operation \n%s", err)
		return
	}
}
