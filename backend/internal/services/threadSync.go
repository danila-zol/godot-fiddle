package services

import "gamehangar/internal/domain/models"

type ForumRepository interface {
	CreateThread(thread models.Thread) (*models.Thread, error)
	UpdateThread(id int, thread models.Thread) (*models.Thread, error)
	DeleteThread(id int) error

	CreateMessage(message models.Message) (*models.Message, error)
	UpdateMessage(id int, message models.Message) (*models.Message, error)
	DeleteMessage(id int) error
}

type ThreadSyncer struct {
	threadRepository ForumRepository
	demoTopicID      int
}

func NewThreadSyncer(r ForumRepository, demoTopicID int) *ThreadSyncer {
	return &ThreadSyncer{
		threadRepository: r,
		demoTopicID:      demoTopicID, // Demo will be stored in a single topic
	}
}

func (s *ThreadSyncer) PostThread(demo models.Demo) (*int, error) {
	thread := models.Thread{
		Title:          demo.Title,
		UserID:         demo.UserID,
		CreatedAt:      demo.CreatedAt,
		LastUpdate:     demo.UpdatedAt,
		TotalUpvotes:   demo.Upvotes,
		TotalDownvotes: demo.Downvotes,
		Tags:           demo.Tags,
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
