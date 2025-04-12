package services

import (
	"gamehangar/internal/domain/models"

	"github.com/google/uuid"
)

type ForumRepository interface {
	CreateThread(thread models.Thread) (*models.Thread, error)
	UpdateThread(id string, thread models.Thread) (*models.Thread, error)
	DeleteThread(id string) error

	CreateMessage(message models.Message) (*models.Message, error)
	UpdateMessage(id string, message models.Message) (*models.Message, error)
	DeleteMessage(id string) error
}

type ThreadSyncer struct {
	threadRepository ForumRepository
	demoTopicID      string
}

func NewThreadSyncer(r ForumRepository, demoTopicID string) *ThreadSyncer {
	return &ThreadSyncer{
		threadRepository: r,
		demoTopicID:      demoTopicID, // Demo will be stored in a single topic
	}
}

func (s *ThreadSyncer) PostThread(demo models.Demo) (*string, error) {
	thread := models.Thread{
		Title:          demo.Title,
		UserID:         demo.UserID,
		CreatedAt:      demo.CreatedAt,
		LastUpdate:     demo.UpdatedAt,
		TotalUpvotes:   demo.Upvotes,
		TotalDownvotes: demo.Downvotes,
		Tags:           demo.Tags,
	}

	if thread.ID == nil {
		id := uuid.NewString()
		thread.ID = &id
	}
	if thread.TopicID == nil {
		id := s.demoTopicID
		thread.TopicID = &id
	}

	t, err := s.threadRepository.CreateThread(thread)
	if err != nil {
		return nil, err
	}

	return t.ID, nil
}

func (s *ThreadSyncer) PatchThread(demo models.Demo) error {
	thread := models.Thread{
		Title:          demo.Title,
		LastUpdate:     demo.UpdatedAt,
		TotalUpvotes:   demo.Upvotes,
		TotalDownvotes: demo.Downvotes,
		Tags:           demo.Tags,
	}

	_, err := s.threadRepository.UpdateThread(*demo.ThreadID, thread)
	if err != nil {
		return err
	}
	return nil
}
