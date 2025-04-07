package repository

import "gamehangar/internal/domain/models"

type UserRepository interface {
	CreateUser(user models.User) error
	GetUsers() ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(id string, user models.User) error
	DeleteUser(id string) error
}

type RoleRepository interface {
	CreateRole(role models.Role) error
	DeleteRole(id string) error
}

type SessionRepository interface {
	CreateSession(session models.Session) error
	DeleteSession(sessionId string) error
}
