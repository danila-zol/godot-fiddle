package repository

import "gamehangar/internal/domain/models"

type DemoRepository interface {
	CreateDemo(demo models.Demo) error
	GetDemos() ([]models.Demo, error)
	GetDemoByID(id string) (models.Demo, error)
	UpdateDemo(id string, demo models.Demo) error
	DeleteDemo(id string) error
}
