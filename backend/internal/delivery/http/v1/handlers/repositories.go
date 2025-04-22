package handlers

import "gamehangar/internal/domain/models"

type AssetRepository interface {
	CreateAsset(asset models.Asset) (*models.Asset, error)
	FindAssets() (*[]models.Asset, error)
	FindAssetByID(id int) (*models.Asset, error)
	UpdateAsset(id int, asset models.Asset) (*models.Asset, error)
	DeleteAsset(id int) error
	NotFoundErr() error
}

type DemoRepository interface {
	CreateDemo(demo models.Demo) (*models.Demo, error)
	FindDemos() (*[]models.Demo, error)
	FindDemoByID(id int) (*models.Demo, error)
	FindDemosByQuery(query string) (*[]models.Demo, error)
	// FindDemosByDate(time time.Time) (*[]models.Demo, error)
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
	UpdateThread(id int, thread models.Thread) (*models.Thread, error)
	DeleteThread(id int) error

	CreateMessage(message models.Message) (*models.Message, error)
	FindMessages() (*[]models.Message, error)
	FindMessagesByThreadID(threadID int) (*[]models.Message, error)
	FindMessageByID(id int) (*models.Message, error)
	UpdateMessage(id int, message models.Message) (*models.Message, error)
	DeleteMessage(id int) error

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
	DeleteAllUserSessions(userID string) error

	NotFoundErr() error
}
