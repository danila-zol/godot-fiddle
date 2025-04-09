package repository

import "gamehangar/internal/domain/models"

type DemoRepository interface {
	CreateDemo(demo models.Demo) error
	FindDemos() ([]models.Demo, error)
	FindDemoByID(id string) (models.Demo, error)
	UpdateDemo(id string, demo models.Demo) error
	DeleteDemo(id string) error
}
