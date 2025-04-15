package services

import "gamehangar/internal/domain/models"

type UserRepository interface {
	FindUserByEmail(email string) (user *models.User, err error)
	FindUserByUsername(username string) (user *models.User, err error)
}

type UserLookup struct {
	userRepository UserRepository
}

func NewUserLookup(r UserRepository) *UserLookup {
	return &UserLookup{
		userRepository: r,
	}
}

func (l *UserLookup) LookupUser(email, username *string) (user *models.User, err error) {
	if email != nil {
		user, err = l.userRepository.FindUserByEmail(*email)
		if err != nil {
			return nil, err
		}
	}
	if username != nil {
		user, err = l.userRepository.FindUserByUsername(*username)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
