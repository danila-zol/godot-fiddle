package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
)

type AssetRoutes struct {
	handler    *handlers.AssetHandler
	authorizer Authorizer
}

func NewAssetRoutes(h *handlers.AssetHandler, a Authorizer) *AssetRoutes {
	return &AssetRoutes{
		handler:    h,
		authorizer: a,
	}
}

func (r *AssetRoutes) InitRoutes(e *echo.Echo) {
	assetGroup := e.Group("/game-hangar/v1/assets")

	protectedAssetGroup := assetGroup.Group("")
	protectedAssetGroup.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		EnforceHandler: r.authorizer.CheckPermissions,
	}))

	protectedAssetGroup.POST("", r.handler.PostAsset)
	assetGroup.GET("/:id", r.handler.GetAssetById)
	assetGroup.GET("", r.handler.GetAssets)
	protectedAssetGroup.PATCH("/:id", r.handler.PatchAsset)
	protectedAssetGroup.DELETE("/:id", r.handler.DeleteAsset)
}
