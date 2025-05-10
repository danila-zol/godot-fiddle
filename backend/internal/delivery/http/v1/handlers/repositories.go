package handlers

import (
	"gamehangar/internal/domain/models"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CreateAsset(asset models.Asset) (*models.Asset, error)
	FindAssets(query []string, limit uint64, order string) (*[]models.Asset, error)
	FindAssetByID(id int) (*models.Asset, error)
	UpdateAsset(id int, asset models.Asset) (*models.Asset, error)
	DeleteAsset(id int) error
	NotFoundErr() error
	ConflictErr() error
}

type DemoRepository interface {
	CreateDemo(demo models.Demo) (*models.Demo, error)
	FindDemos(query []string, limit uint64, order string) (*[]models.Demo, error)
	FindDemoByID(id int) (*models.Demo, error)
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
	FindThreads(query []string, limit uint64, order string) (*[]models.Thread, error)
	FindThreadByID(id int) (*models.Thread, error)
	UpdateThread(id int, thread models.Thread) (*models.Thread, error)
	DeleteThread(id int) error

	CreateMessage(message models.Message) (*models.Message, error)
	FindMessages(query []string, limit uint64, order string) (*[]models.Message, error)
	FindMessagesByThreadID(threadID int) (*[]models.Message, error)
	FindMessageByID(id int) (*models.Message, error)
	UpdateMessage(id int, message models.Message) (*models.Message, error)
	DeleteMessage(id int) error

	NotFoundErr() error
	ConflictErr() error
}

type UserRepository interface {
	CreateUser(user models.User) (*models.User, error)
	FindUsers(query []string, limit uint64) (*[]models.User, error)
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
