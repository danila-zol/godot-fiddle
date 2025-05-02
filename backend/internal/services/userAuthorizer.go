package services

import (
	"gamehangar/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthorizerRepository interface {
	FindUserByID(id string) (*models.User, error)
	FindUserByEmail(email string) (user *models.User, err error)
	FindUserByUsername(username string) (user *models.User, err error)
	NotFoundErr() error
}

type UserAuthorizer struct {
	repository UserAuthorizerRepository
}

func NewUserAuthorizer(r UserAuthorizerRepository) *UserAuthorizer {
	return &UserAuthorizer{
		repository: r,
	}
}

func (a *UserAuthorizer) IdentifyUser(email, username *string) (user *models.User, err error) {
	if email != nil {
		user, err = a.repository.FindUserByEmail(*email)
		if err == nil {
			return user, nil
		}
	}
	if username != nil {
		user, err = a.repository.FindUserByUsername(*username)
		if err == nil {
			return user, nil
		}
	}

	return nil, a.repository.NotFoundErr()
}

func (a *UserAuthorizer) CreatePasswordHash(password *string) (*string, error) {
	var hash string

	h, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hash = string(h)

	return &hash, nil
}

func (a *UserAuthorizer) CheckPassword(password, userID *string) error {
	user, err := a.repository.FindUserByID(*userID)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password))
}
