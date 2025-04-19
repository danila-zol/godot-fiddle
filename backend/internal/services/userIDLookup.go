package services

import "gamehangar/internal/domain/models"

type UsernameResolver interface {
	FindUserByEmail(email string) (user *models.User, err error)
	FindUserByUsername(username string) (user *models.User, err error)
	NotFoundErr() error
}

type UserIdentifier struct {
	usernameResolver UsernameResolver
}

func NewUserIdentifier(r UsernameResolver) *UserIdentifier {
	return &UserIdentifier{
		usernameResolver: r,
	}
}

func (i *UserIdentifier) IdentifyUser(email, username *string) (user *models.User, err error) {
	if email != nil {
		user, err = i.usernameResolver.FindUserByEmail(*email)
		if err == nil {
			return user, nil
		}
	}
	if username != nil {
		user, err = i.usernameResolver.FindUserByUsername(*username)
		if err == nil {
			return user, nil
		}
	}

	return nil, i.usernameResolver.NotFoundErr()
}
