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

type DemoRepository interface {
	FindDemoByID(id int) (*models.Demo, error)
	NotFoundErr() error
}

type ThreadSyncer struct {
	threadRepository ForumRepository
	demoRepository   DemoRepository
	demoTopicID      int
}

func NewThreadSyncer(fr ForumRepository, dr DemoRepository, demoTopicID int) *ThreadSyncer {
	return &ThreadSyncer{
		threadRepository: fr,
		demoRepository:   dr,
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

func (s *ThreadSyncer) PatchThread(demoID int, demo models.Demo) error {
	thread := models.Thread{
		Title:          demo.Title,
		LastUpdate:     demo.UpdatedAt,
		TotalUpvotes:   demo.Upvotes,
		TotalDownvotes: demo.Downvotes,
		Tags:           demo.Tags,
	}

	if demo.ThreadID == nil {
		d, err := s.demoRepository.FindDemoByID(demoID)
		if notFoundErr := s.demoRepository.NotFoundErr(); notFoundErr == err {
			return notFoundErr
		}
		if err != nil {
			return err
		}
		demo.ThreadID = d.ThreadID
	}
	_, err := s.threadRepository.UpdateThread(*demo.ThreadID, thread)
	if err != nil {
		return err
	}
	return nil
}
