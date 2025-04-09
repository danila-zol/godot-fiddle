package repository

import "gamehangar/internal/domain/models"

type TopicRepository interface {
	CreateTopic(topic models.Topic) error
	FindTopics() (*[]models.Topic, error)
	FindTopicByID(id string) (*models.Topic, error)
	UpdateTopic(id string, topic models.Topic) error
	DeleteTopic(id string) error
}

type ThreadRepository interface {
	CreateThread(thread models.Thread) error
	FindThreads() (*[]models.Thread, error)
	FindThreadByID(id string) (*models.Thread, error)
	UpdateThread(id string, thread models.Thread) error
	DeleteThread(id string) error
}

type MessageRepository interface {
	CreateMessage(message models.Message) error
	FindMessages() (*[]models.Message, error)
	FindMessagesByThreadID(threadID string) (*[]models.Message, error)
	FindMessageByID(id string) (*models.Message, error)
	UpdateMessage(id string, message models.Message) error
	DeleteMessage(id string) error
}
