package psqlRepository

import (
	// "context"
	"errors"
	// "gamehangar/internal/config/psqlDatabseConfig"
	// "gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	// "gamehangar/pkg/ternMigrate"
	// "os"
	"testing"
	"time"

	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// independent	bool = false
	// testDBClient     *psqlDatabase.PsqlDatabaseClient
	// testEnforcer *psqlCasbinClient.CasbinClient

	// roleID          uuid.UUID
	// userID          uuid.UUID

	topicID          int          = 1
	topicName        string       = "Test"
	topicNameUpdated string       = "Test UPDATE"
	topicVersion     int          = 1
	topic            models.Topic = models.Topic{Name: &topicName}
	topicUpdated     models.Topic = models.Topic{Name: &topicNameUpdated, Version: &topicVersion}

	threadID           int           = 1
	threadTitle        string        = "Test"
	threadTitleUpdated string        = "Test UPDATE"
	threadTags         []string      = []string{"TEST", "test", "thread"}
	thread             models.Thread = models.Thread{Title: &threadTitle, Tags: &threadTags, UserID: &userID, TopicID: &topicID}
	threadUpdated      models.Thread = models.Thread{Title: &threadTitleUpdated}

	messageID            int            = 1
	messageTitle         string         = "Test"
	messsageTitleUpdated string         = "Test UPDATE"
	messageBody          string         = "An demo for integration testing for PSQL Repo"
	messageTags          []string       = []string{"TEST", "test", "message"}
	message              models.Message = models.Message{Title: &messageTitle, Body: &messageBody, UserID: &userID, ThreadID: &threadID, Tags: &messageTags}
	messageUpdated       models.Message = models.Message{Title: &messsageTitleUpdated}
)

func init() {
	if independent {
		ResetDB()
	}
}

func TestCreateTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	_, err := r.CreateTopic(topic)
	assert.NoError(t, err)
}

func TestFindTopicByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopicByID(topicID)
	assert.NoError(t, err)
}

func TestFindTopicByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopicByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindTopics(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopics()
	assert.NoError(t, err)
}

func TestUpdateTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	topicID = 1
	resultTopic, err := r.UpdateTopic(topicID, topicUpdated)
	assert.NoError(t, err)

	modifiedTopic := topic
	modifiedTopic.ID = &topicID
	modifiedTopic.Name = &topicNameUpdated
	newVersion := *topicUpdated.Version + 1
	modifiedTopic.Version = &newVersion

	assert.Equal(t, modifiedTopic, *resultTopic)
}

func TestUpdateTopicMultiple(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	topicID = 1

	modifiedTopic := topic // Manual update
	modifiedTopic.ID = &topicID

	for i := 2; i < 6; i++ {
		newName := *topicUpdated.Name + " New"
		topicUpdated.Name = &newName
		newVersion := i
		topicUpdated.Version = &newVersion

		resultTopic, err := r.UpdateTopic(topicID, topicUpdated)
		assert.NoError(t, err)

		newerVersion := i + 1
		modifiedTopic.Version = &newerVersion
		modifiedTopic.Name = &newName

		assert.Equal(t, modifiedTopic, *resultTopic)
	}
}

func TestUpdateTopicConflict(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	topicID = 1
	_, err := r.UpdateTopic(topicID, topicUpdated)
	if assert.Error(t, err) {
		assert.Equal(t, r.conflictErr, err)
	}
}

func TestDeleteTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer, conflictErr: errors.New("Record conflict!")}
	err := r.DeleteTopic(topicID)
	assert.NoError(t, err)
}

func TestCreateThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	topic, err := r.CreateTopic(topic)
	topicID = *topic.ID
	th, err := r.CreateThread(thread)
	threadID = *th.ID
	assert.NoError(t, err)
}

func TestFindThreadByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	thread, err := r.FindThreadByID(threadID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *thread.Views)
	}
}

func TestFindThreadByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindThreadByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindThreads(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindThreads(nil, 0, "")
	assert.NoError(t, err)
}

func TestFindThreadsByQuery(t *testing.T) {
	var (
		threadTitleAlt string        = "Cheeseboiger"
		threadTagsAlt  []string      = []string{"The Magnificent Seven", "Rock the Casbah"}
		threadAlt      models.Thread = models.Thread{
			Title:   &threadTitleAlt,
			Tags:    &threadTagsAlt,
			UserID:  &userID,
			TopicID: &topicID,
		}

		threadTitleAltRu string        = "Стук"
		threadTagsAltRu  []string      = []string{"Cheeseboiger", "Муравейник"}
		threadAltRu      models.Thread = models.Thread{
			Title:   &threadTitleAltRu,
			Tags:    &threadTagsAltRu,
			UserID:  &userID,
			TopicID: &topicID,
		}
	)

	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}

	for q, th := range map[string]models.Thread{"The Magnificent Seven": threadAlt, "стук": threadAltRu} {
		resultThread, err := r.CreateThread(th)
		assert.NoError(t, err)

		queryThreads, err := r.FindThreads([]string{q}, 0, "highest-rated")
		t.Log(queryThreads)
		if assert.NoError(t, err) {
			queriedThread := *queryThreads
			assert.Equal(t, resultThread.Title, queriedThread[0].Title)
		}
	}

	// Try to query both and check ordering
	threads, err := r.FindThreads([]string{"cheeseboiger"}, 0, "newest-updated")
	if assert.NoError(t, err) {
		th := *threads
		assert.Len(t, th, 2)
		var timeOrder, timeOrderExpected []time.Time
		timeOrderExpected = []time.Time{*th[0].UpdatedAt, *th[1].UpdatedAt}
		for _, m := range th {
			timeOrder = append(timeOrder, *m.UpdatedAt)
		}
		assert.Equal(
			t,
			timeOrderExpected,
			timeOrder,
		)
	}
	// Query with limit
	threads, err = r.FindThreads([]string{"cheeseboiger"}, 1, "most-views")
	if assert.NoError(t, err) {
		assert.Len(t, *threads, 1)
	}
}

func TestUpdateThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}

	oldThread, err := r.FindThreadByID(threadID)
	assert.NoError(t, err)

	resultThread, err := r.UpdateThread(threadID, threadUpdated)
	if assert.NoError(t, err) {
		assert.Equal(t, oldThread.CreatedAt, resultThread.CreatedAt)
		assert.Equal(t, oldThread.Tags, resultThread.Tags)
		assert.Equal(t, oldThread.Rating, resultThread.Rating)
		assert.Equal(t, oldThread.Views, resultThread.Views)

		assert.NotEqual(t, oldThread.UpdatedAt, resultThread.UpdatedAt)

		assert.NotEqual(t, oldThread.Title, resultThread.Title)
		assert.Equal(t, threadUpdated.Title, resultThread.Title)
	}
}

func TestDeleteThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	err := r.DeleteThread(threadID)
	assert.NoError(t, err)
}

func TestCreateMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	thread, err := r.CreateThread(thread)
	threadID = *thread.ID
	_, err = r.CreateMessage(message)
	assert.NoError(t, err)
}

func TestFindMessageByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	message, err := r.FindMessageByID(messageID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *message.Views)
	}
}

func TestFindMessageByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindMessageByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindMessageByThreadID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindMessagesByThreadID(threadID)
	assert.NoError(t, err)
}

func TestFindMessageByThreadIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindMessagesByThreadID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindMessagesByQuery(t *testing.T) {
	var (
		messageTitleAlt string         = "The Magnificent Seven"
		messageBodyAlt  string         = "Marx was skint but he had sense, Engels lent him the necessary pence"
		messageTagsAlt  []string       = []string{"Cheeseboiger", "Rock the Casbah"}
		messageAlt      models.Message = models.Message{Title: &messageTitleAlt, Body: &messageBodyAlt, Tags: &messageTagsAlt, ThreadID: &threadID, UserID: &userID}

		messageTitleAltRu string         = "Стук"
		messageBodyAltRu  string         = `Я скажу одно лишь слово: "Cheeseboiger"`
		messageAltRu      models.Message = models.Message{Title: &messageTitleAltRu, Body: &messageBodyAltRu, ThreadID: &threadID, UserID: &userID}
	)

	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}

	for q, m := range map[string]models.Message{"seven": messageAlt, "стук": messageAltRu} {
		resultMessage, err := r.CreateMessage(m)
		assert.NoError(t, err)

		queryMessages, err := r.FindMessages([]string{q}, 0, "highest-rated")
		if assert.NoError(t, err) {
			queriedMessage := *queryMessages
			assert.Equal(t, resultMessage.Title, queriedMessage[0].Title)
		}
	}

	// Try to query both and check ordering
	messages, err := r.FindMessages([]string{"cheeseboiger"}, 0, "newest-updated")
	if assert.NoError(t, err) {
		m := *messages
		assert.Len(t, m, 2)
		var timeOrder, timeOrderExpected []time.Time
		timeOrderExpected = []time.Time{*m[0].UpdatedAt, *m[1].UpdatedAt}
		for _, m := range m {
			timeOrder = append(timeOrder, *m.UpdatedAt)
		}
		assert.Equal(
			t,
			timeOrderExpected,
			timeOrder,
		)
	}
	// Query with limit
	messages, err = r.FindMessages([]string{"cheeseboiger"}, 1, "most-views")
	if assert.NoError(t, err) {
		assert.Len(t, *messages, 1)
	}
}

func TestFindMessages(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindMessages(nil, 0, "")
	assert.NoError(t, err)
}

func TestUpdateMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}

	oldMessasge, err := r.FindMessageByID(messageID)
	assert.NoError(t, err)

	resultMessage, err := r.UpdateMessage(messageID, messageUpdated)
	if assert.NoError(t, err) {
		assert.Equal(t, oldMessasge.CreatedAt, resultMessage.CreatedAt)
		assert.Equal(t, oldMessasge.Tags, resultMessage.Tags)
		assert.Equal(t, oldMessasge.Rating, resultMessage.Rating)
		assert.Equal(t, oldMessasge.Views, resultMessage.Views)

		assert.NotEqual(t, oldMessasge.UpdatedAt, resultMessage.UpdatedAt)

		assert.NotEqual(t, oldMessasge.Title, resultMessage.Title)
		assert.Equal(t, messageUpdated.Title, resultMessage.Title)
	}
}

func TestDeleteMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	err := r.DeleteMessage(messageID)
	if assert.NoError(t, err) {
		teardownForum(&r)
	}
}

func teardownForum(r *PsqlForumRepository) {
	remainderTopics, err := r.FindTopics()
	if err != nil {
		panic(err)
	}
	for _, t := range *remainderTopics {
		err = r.DeleteTopic(*t.ID)
		if err != nil {
			panic(err)
		}
	}
}
