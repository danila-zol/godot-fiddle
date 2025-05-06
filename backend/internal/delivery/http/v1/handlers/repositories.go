package handlers

import (
	"gamehangar/internal/domain/models"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CreateAsset(asset models.Asset) (*models.Asset, error)
	FindAssets() (*[]models.Asset, error)
	FindAssetByID(id int) (*models.Asset, error)
	FindAssetsByQuery(query *[]string) (*[]models.Asset, error)
	UpdateAsset(id int, asset models.Asset) (*models.Asset, error)
	DeleteAsset(id int) error
	NotFoundErr() error
	ConflictErr() error
}

type DemoRepository interface {
	CreateDemo(demo models.Demo) (*models.Demo, error)
	FindDemos() (*[]models.Demo, error)
	FindDemoByID(id int) (*models.Demo, error)
	FindDemosByQuery(query *[]string) (*[]models.Demo, error)
	UpdateDemo(id int, demo models.Demo) (*models.Demo, error)
	DeleteDemo(id int) error
	NotFoundErr() error
}

type ForumRepository interface {
	CreateTopic(topic models.Topic) (*models.Topic, error)
	FindTopics() (*[]models.Topic, error)
	FindTopicByID(id int) (*models.Topic, error)
	UpdateTopic(id int, topic models.Topic) (*models.Topic, error)
	DeleteTopic(id int) error

	CreateThread(thread models.Thread) (*models.Thread, error)
	FindThreads() (*[]models.Thread, error)
	FindThreadByID(id int) (*models.Thread, error)
	FindThreadsByQuery(query *[]string) (*[]models.Thread, error)
	UpdateThread(id int, thread models.Thread) (*models.Thread, error)
	DeleteThread(id int) error

	CreateMessage(message models.Message) (*models.Message, error)
	FindMessages() (*[]models.Message, error)
	FindMessagesByQuery(query *[]string) (*[]models.Message, error)
	FindMessagesByThreadID(threadID int) (*[]models.Message, error)
	FindMessageByID(id int) (*models.Message, error)
	UpdateMessage(id int, message models.Message) (*models.Message, error)
	DeleteMessage(id int) error

	NotFoundErr() error
	ConflictErr() error
}

type UserRepository interface {
	CreateUser(user models.User) (*models.User, error)
	FindUsers() (*[]models.User, error)
	FindUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, user models.User) (*models.User, error)
	DeleteUser(id uuid.UUID) error

	CreateRole(role models.Role) (*models.Role, error)
	FindRoleByID(id uuid.UUID) (*models.Role, error)
	UpdateRole(id uuid.UUID, role models.Role) (*models.Role, error)
	DeleteRole(id uuid.UUID) error

	CreateSession(session models.Session) (*models.Session, error)
	FindSessionByID(id uuid.UUID) (*models.Session, error)
	DeleteSession(id uuid.UUID) error
	DeleteAllUserSessions(userid uuid.UUID) error

	NotFoundErr() error
	ConflictErr() error
}
