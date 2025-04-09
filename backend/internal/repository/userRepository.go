package repository

import "gamehangar/internal/domain/models"

type UserRepository interface {
	CreateUser(user models.User) error
	FindUsers() ([]models.User, error)
	FindUserByID(id string) (models.User, error)
	UpdateUser(id string, user models.User) error
	DeleteUser(id string) error
}

type RoleRepository interface {
	CreateRole(role models.Role) error
	FindRoleByID(id string) (models.Role, error)
	UpdateRole(id string, role models.Role) error
	DeleteRole(id string) error
}

type SessionRepository interface {
	CreateSession(session models.Session) error
	FindSessionByID(id string) (models.Session, error)
	DeleteSession(id string) error
}
