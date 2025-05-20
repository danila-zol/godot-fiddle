package services

import (
	"gamehangar/internal/domain/models"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthorizerRepository interface {
	FindSessionByID(id uuid.UUID) (*models.Session, error)
	FindUserByID(id uuid.UUID) (*models.User, error)
	FindUserByEmail(email string) (user *models.User, err error)
	FindUserByUsername(username string) (user *models.User, err error)
	NotFoundErr() error
}

type Enforcer interface {
	Enforce(sub, obj, act string) (bool, error)
}

type UserAuthorizer struct {
	repository UserAuthorizerRepository
	enforcer   Enforcer
}

func NewUserAuthorizer(r UserAuthorizerRepository, e Enforcer) *UserAuthorizer {
	return &UserAuthorizer{
		repository: r,
		enforcer:   e,
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

func (a *UserAuthorizer) CheckPassword(password *string, userID uuid.UUID) error {
	user, err := a.repository.FindUserByID(userID)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password))
}

func (a *UserAuthorizer) CheckPermissions(c echo.Context, user string) (bool, error) {
	var (
		obj, act  string
		sessionID uuid.UUID
	)

	cookie, err := c.Cookie("sessionID")
	if err != nil {
		sessionSlice, ok := c.Request().Header["Sessionid"]
		if !ok {
			return false, err
		}
		sessionID, err = uuid.Parse(sessionSlice[0])
	} else {
		sessionID, err = uuid.Parse(cookie.Value)
	}
	if err != nil {
		return false, err
	}

	session, err := a.repository.FindSessionByID(sessionID)
	sub, err := a.repository.FindUserByID(*session.UserID)
	if err != nil {
		return false, err
	}
	c.Set("userTier", *sub.Role)

	obj = strings.TrimPrefix(c.Request().URL.Path, "/game-hangar/v1/")

	act = c.Request().Method

	eft1, err := a.enforcer.Enforce(sub.ID.String(), obj, act) // Check user permissions over the obj
	eft2, err := a.enforcer.Enforce(*sub.Role, obj, act)       // Check user role permissions over the obj

	return (eft1 || eft2), err
}
