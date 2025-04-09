package repository

import "gamehangar/internal/domain/models"

type AssetRepository interface {
	CreateAsset(asset models.Asset) error
	FindAssets() (*[]models.Asset, error)
	FindAssetByID(id string) (*models.Asset, error)
	UpdateAsset(id string, asset models.Asset) error
	DeleteAsset(id string) error
}
