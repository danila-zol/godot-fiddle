package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	"github.com/labstack/echo/v4"
)

type AssetRoutes struct {
	handler *handlers.AssetHandler
}

func NewAssetRoutes(h *handlers.AssetHandler) *AssetRoutes {
	return &AssetRoutes{
		handler: h,
	}
}

func (r *AssetRoutes) InitRoutes(e *echo.Echo) {
	assetGroup := e.Group("/v1/asset")

	protectedAssetGroup := assetGroup.Group("/protected")

	protectedAssetGroup.POST("", r.handler.PostAsset)
	assetGroup.GET("/:id", r.handler.GetAssetById)
	assetGroup.GET("", r.handler.GetAssets)
	protectedAssetGroup.PATCH("/:id", r.handler.PatchAsset)
	protectedAssetGroup.DELETE("/:id", r.handler.DeleteAsset)
}
