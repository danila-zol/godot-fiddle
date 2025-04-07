package repository

import "gamehangar/internal/domain/models"

type TopicRepository interface {
	CreateTopic(topic models.Topic) error
	GetTopics() ([]models.Topic, error)
	GetTopicByID(id string) (models.Topic, error)
	UpdateTopic(id string, topic models.Topic) error
	DeleteTopic(id string) error
}

type ThreadRepository interface {
	CreateThread(thread models.Thread) error
	GetThreads() ([]models.Thread, error)
	GetThreadByID(id string) (models.Thread, error)
	UpdateThread(id string, thread models.Thread) error
	DeleteThread(id string) error
}

type MessageRepository interface {
	CreateMessage(message models.Message) error
	GetMessages() ([]models.Message, error)
	GetMessagesByThreadId(threadId string) ([]models.Message, error)
	GetMessageByID(id string) (models.Message, error)
	UpdateMessage(id string, message models.Message) error
	DeleteMessage(id string) error
}
