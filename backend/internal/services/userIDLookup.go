package services

import "gamehangar/internal/domain/models"

type UserRepository interface {
	FindUserByEmail(email string) (user *models.User, err error)
	FindUserByUsername(username string) (user *models.User, err error)
}

type UserIdentifier struct {
	userRepository UserRepository
}

func NewUserIdentifier(r UserRepository) *UserIdentifier {
	return &UserIdentifier{
		userRepository: r,
	}
}

func (l *UserIdentifier) IdentifyUser(email, username *string) (user *models.User, err error) {
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
