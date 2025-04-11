package handlers

import "gamehangar/internal/domain/models"

type AssetRepository interface {
	CreateAsset(asset models.Asset) (*models.Asset, error)
	FindAssets() (*[]models.Asset, error)
	FindAssetByID(id string) (*models.Asset, error)
	UpdateAsset(id string, asset models.Asset) (*models.Asset, error)
	DeleteAsset(id string) error
	NotFoundErr() error
}

type DemoRepository interface {
	CreateDemo(demo models.Demo) (*models.Demo, error)
	FindDemos() (*[]models.Demo, error)
	FindDemoByID(id string) (*models.Demo, error)
	UpdateDemo(id string, demo models.Demo) (*models.Demo, error)
	DeleteDemo(id string) error
	NotFoundErr() error
}

type ForumRepository interface {
	CreateTopic(topic models.Topic) (*models.Topic, error)
	FindTopics() (*[]models.Topic, error)
	FindTopicByID(id string) (*models.Topic, error)
	UpdateTopic(id string, topic models.Topic) (*models.Topic, error)
	DeleteTopic(id string) error

	CreateThread(thread models.Thread) (*models.Thread, error)
	FindThreads() (*[]models.Thread, error)
	FindThreadByID(id string) (*models.Thread, error)
	UpdateThread(id string, thread models.Thread) (*models.Thread, error)
	DeleteThread(id string) error

	CreateMessage(message models.Message) (*models.Message, error)
	FindMessages() (*[]models.Message, error)
	FindMessagesByThreadID(threadID string) (*[]models.Message, error)
	FindMessageByID(id string) (*models.Message, error)
	UpdateMessage(id string, message models.Message) (*models.Message, error)
	DeleteMessage(id string) error

	NotFoundErr() error
}

type UserRepository interface {
	CreateUser(user models.User) (*models.User, error)
	FindUsers() (*[]models.User, error)
	FindUserByID(id string) (*models.User, error)
	UpdateUser(id string, user models.User) (*models.User, error)
	DeleteUser(id string) error

	CreateRole(role models.Role) (*models.Role, error)
	FindRoleByID(id string) (*models.Role, error)
	UpdateRole(id string, role models.Role) (*models.Role, error)
	DeleteRole(id string) error

	CreateSession(session models.Session) (*models.Session, error)
	FindSessionByID(id string) (*models.Session, error)
	DeleteSession(id string) error

	NotFoundErr() error
}
