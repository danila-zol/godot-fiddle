package services

import (
	"gamehangar/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindUserByID(id string) (*models.User, error)
}

type PasswordManager struct {
	userRepository UserRepository
}

func NewPasswordManager(r UserRepository) *PasswordManager {
	return &PasswordManager{
		userRepository: r,
	}
}

func (m *PasswordManager) CreatePasswordHash(password *string) (*string, error) {
	var hash string

	h, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hash = string(h)

	return &hash, nil
}

func (m *PasswordManager) CheckPassword(password, userID *string) error {
	user, err := m.userRepository.FindUserByID(*userID)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password))
}
