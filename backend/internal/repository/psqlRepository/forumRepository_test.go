package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"gamehangar/pkg/ternMigrate"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// testDBClient     *psqlDatabase.PsqlDatabaseClient

	// roleID          uuid.UUID
	// userID          uuid.UUID

	topicID          int
	topicName        string       = "Test"
	topicNameUpdated string       = "Test UPDATE"
	topicVersion     int          = 1
	topic            models.Topic = models.Topic{Name: &topicName}
	topicUpdated     models.Topic = models.Topic{Name: &topicNameUpdated, Version: &topicVersion}

	threadID           int
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
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../../.env")
	if err != nil {
		panic("Error loading .env file:" + err.Error() + ": " + wd)
	}
	databaseConfig, err := psqlDatabseConfig.PsqlConfig{}.NewConfig(
		psqlDatabase.MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"),
	)
	if err != nil {
		panic("Error loading PSQL database Config")
	}
	testDBClient, err = psqlDatabase.PsqlDatabase{}.NewDatabaseClient(
		os.Getenv("PSQL_CONNSTRING"), ternMigrate.Migrator{}, databaseConfig,
	)
	if err != nil {
		panic("Error setting up new DatabaseClient")
	}
	c, _ := testDBClient.AcquireConn() // WARNING! Integration tests DROP TABLEs
	_, err = c.Exec(context.Background(), `
		DROP TRIGGER IF EXISTS increment_topic_version_on_update ON forum.topics; 
		DROP SCHEMA IF EXISTS "user" CASCADE;
		DROP SCHEMA IF EXISTS "forum" CASCADE;

		DROP COLLATION IF EXISTS case_insensitive;
		DROP INDEX IF EXISTS thread_gin_index_ts;
		DROP INDEX IF EXISTS message_gin_index_ts;

		CREATE SCHEMA IF NOT EXISTS "user";

		CREATE TABLE "user".roles (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"name" VARCHAR(255) NOT NULL
		-- "permissions" VARCHAR(64)[]
		);

		CREATE TABLE "user".users (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"username" VARCHAR(255) NOT NULL UNIQUE,
		"display_name" VARCHAR(255),
		"email" VARCHAR(255) NOT NULL UNIQUE,
		"password" VARCHAR(255) NOT NULL,
		"verified" BOOLEAN NOT NULL DEFAULT false,
		"role_id" UUID NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"karma" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE "user".sessions (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"user_id" UUID NOT NULL REFERENCES "user".users (id) ON DELETE CASCADE
		);

		CREATE SCHEMA IF NOT EXISTS forum;

		CREATE TABLE forum.topics (
		"id"  INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"name" VARCHAR(255) NOT NULL,
		"version" INTEGER NOT NULL DEFAULT 1
		);

		CREATE TABLE forum.threads (
		"id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"title" VARCHAR(255) NOT NULL,
		"user_id" UUID NOT NULL,
		"topic_id" INTEGER NOT NULL REFERENCES forum.topics (id) ON DELETE CASCADE,
		"tags" VARCHAR(255)[],
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"upvotes" INTEGER NOT NULL DEFAULT 0,
		"downvotes" INTEGER NOT NULL DEFAULT 0,
		"views" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE forum.messages (
		"id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"thread_id" INTEGER NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
		"user_id" UUID NOT NULL,
		"title" VARCHAR(255) NOT NULL,
		"body" VARCHAR NOT NULL,
		"tags" VARCHAR(255)[],
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"upvotes" INTEGER NOT NULL DEFAULT 0,
		"downvotes" INTEGER NOT NULL DEFAULT 0,
		"views" INTEGER NOT NULL DEFAULT 0
		);

		CREATE OR REPLACE FUNCTION increment_version()
		RETURNS TRIGGER AS
		$func$
		BEGIN
		NEW.version := OLD.version + 1;
		RETURN NEW;
		END;
		$func$ LANGUAGE plpgsql;

		CREATE TRIGGER increment_topic_version_on_update
		BEFORE UPDATE ON forum.topics
		FOR EACH ROW
		EXECUTE FUNCTION increment_version();

		CREATE COLLATION IF NOT EXISTS case_insensitive (provider = icu, locale = 'und-u-ks-level2', deterministic = false);

		CREATE OR REPLACE FUNCTION to_tsvector_multilang(text) RETURNS tsvector AS $$
		BEGIN
		RETURN 
		to_tsvector('english', $1) || 
		to_tsvector('russian', $1);
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;

		CREATE OR REPLACE FUNCTION to_tsquery_multilang(text) RETURNS tsquery AS $$
		BEGIN
		RETURN
		websearch_to_tsquery('english', $1) || 
		websearch_to_tsquery('russian', $1);
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;

		ALTER TABLE forum.threads ADD COLUMN thread_ts tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector_multilang("title"), 'A')
		) STORED;
		CREATE INDEX thread_gin_index_ts ON forum.threads USING GIN (thread_ts);

		ALTER TABLE forum.messages ADD COLUMN message_ts tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector_multilang("title"), 'A') ||
		setweight(to_tsvector_multilang(COALESCE("body", '')), 'B')
		) STORED;
		CREATE INDEX message_gin_index_ts ON forum.messages USING GIN (message_ts);
		`)
	if err != nil {
		panic("Error resetting forum schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".roles (name) VALUES ($1) RETURNING (id)`,
		"admin").Scan(&roleID)
	if err != nil {
		panic("Error resetting forum schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".users
		(username, display_name, email, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING (id)`,
		"mike-pech", "Mike", "test@email.com", "TestPassword", roleID).Scan(&userID)
	if err != nil {
		panic("Error resetting forum schema" + err.Error())
	}
}

func TestCreateTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.CreateTopic(topic)
	assert.NoError(t, err)
}

func TestFindTopicByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopicByID(topicID)
	assert.NoError(t, err)
}

func TestFindTopicByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopicByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindTopics(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindTopics()
	assert.NoError(t, err)
}

func TestUpdateTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
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
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
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
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	topicID = 1
	_, err := r.UpdateTopic(topicID, topicUpdated)
	if assert.Error(t, err) {
		assert.Equal(t, r.conflictErr, err)
	}
}

func TestDeleteTopic(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	err := r.DeleteTopic(topicID)
	assert.NoError(t, err)
}

func TestCreateThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	topic, err := r.CreateTopic(topic)
	topicID = *topic.ID
	_, err = r.CreateThread(thread)
	assert.NoError(t, err)
}

func TestFindThreadByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	thread, err := r.FindThreadByID(threadID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *thread.Views)
	}
}

func TestFindThreadByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	_, err := r.FindThreadByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindThreads(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	_, err := r.FindThreads(nil, 0)
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

	r := PsqlForumRepository{databaseClient: testDBClient}

	for q, th := range map[string]models.Thread{"The Magnificent Seven": threadAlt, "стук": threadAltRu} {
		resultThread, err := r.CreateThread(th)
		assert.NoError(t, err)

		queryThreads, err := r.FindThreads([]string{q}, 0)
		t.Log(queryThreads)
		if assert.NoError(t, err) {
			queriedThread := *queryThreads
			assert.Equal(t, resultThread.Title, queriedThread[0].Title)
		}
	}

	// Try to query both and check ordering
	threads, err := r.FindThreads([]string{"cheeseboiger"}, 0)
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
	threads, err = r.FindThreads([]string{"cheeseboiger"}, 1)
	if assert.NoError(t, err) {
		assert.Len(t, *threads, 1)
	}
}

func TestUpdateThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	resultThread, err := r.UpdateThread(threadID, threadUpdated)
	assert.NoError(t, err)

	modifiedThread := thread
	modifiedThread.ID = &assetID
	modifiedThread.Title = &threadTitleUpdated
	modifiedThread.CreatedAt = resultThread.CreatedAt // Timestamps are created on DB
	modifiedThread.UpdatedAt = resultThread.UpdatedAt
	modifiedThread.Upvotes = resultThread.Upvotes
	modifiedThread.Downvotes = resultThread.Downvotes
	modifiedThread.Views = resultThread.Views

	assert.Equal(t, modifiedThread, *resultThread)
}

func TestDeleteThread(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	err := r.DeleteThread(threadID)
	assert.NoError(t, err)
}

func TestCreateMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	thread, err := r.CreateThread(thread)
	threadID = *thread.ID
	_, err = r.CreateMessage(message)
	assert.NoError(t, err)
}

func TestFindMessageByID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	message, err := r.FindMessageByID(messageID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *message.Views)
	}
}

func TestFindMessageByIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	_, err := r.FindMessageByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindMessageByThreadID(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	_, err := r.FindMessagesByThreadID(threadID)
	assert.NoError(t, err)
}

func TestFindMessageByThreadIDNoRows(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
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

	r := PsqlForumRepository{databaseClient: testDBClient}

	for q, m := range map[string]models.Message{"seven": messageAlt, "стук": messageAltRu} {
		resultMessage, err := r.CreateMessage(m)
		assert.NoError(t, err)

		queryMessages, err := r.FindMessages([]string{q}, 0)
		if assert.NoError(t, err) {
			queriedMessage := *queryMessages
			assert.Equal(t, resultMessage.Title, queriedMessage[0].Title)
		}
	}

	// Try to query both and check ordering
	messages, err := r.FindMessages([]string{"cheeseboiger"}, 0)
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
	messages, err = r.FindMessages([]string{"cheeseboiger"}, 1)
	if assert.NoError(t, err) {
		assert.Len(t, *messages, 1)
	}
}

func TestFindMessages(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	_, err := r.FindMessages(nil, 0)
	assert.NoError(t, err)
}

func TestUpdateMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
	resultMessage, err := r.UpdateMessage(messageID, messageUpdated)
	assert.NoError(t, err)

	modifiedMessage := message
	modifiedMessage.ID = &assetID
	modifiedMessage.Title = &messsageTitleUpdated
	modifiedMessage.CreatedAt = resultMessage.CreatedAt // Timestamps are created on DB
	modifiedMessage.UpdatedAt = resultMessage.UpdatedAt
	modifiedMessage.Upvotes = resultMessage.Upvotes
	modifiedMessage.Downvotes = resultMessage.Downvotes
	modifiedMessage.Views = resultMessage.Views

	assert.Equal(t, modifiedMessage, *resultMessage)
}

func TestDeleteMessage(t *testing.T) {
	r := PsqlForumRepository{databaseClient: testDBClient}
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
