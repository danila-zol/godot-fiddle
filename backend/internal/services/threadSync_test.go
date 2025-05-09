package services

import (
	"context"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"gamehangar/internal/repository/psqlRepository"
	"gamehangar/pkg/ternMigrate"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	threadSyncer    *ThreadSyncer
	forumRepository *psqlRepository.PsqlForumRepository
	demoRepository  *psqlRepository.PsqlDemoRepository

	topicID int
	// roleID string
	// userID string

	demoTitle        string      = "Test Demo"
	demoTitleUpdated string      = "Test UPDATE Demo"
	demoDescription  string      = "An demo for integration testing for PSQL Repo"
	demoLink         string      = "https://example.com"
	demoTags         []string    = []string{"TEST", "test"}
	demo             models.Demo = models.Demo{Title: &demoTitle, Description: &demoDescription, Link: &demoLink, Tags: &demoTags}
	demoUpdated      models.Demo = models.Demo{Title: &demoTitleUpdated}
)

func init() {
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../.env")
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
		DROP SCHEMA IF EXISTS "user" CASCADE;
		DROP SCHEMA IF EXISTS "demo" CASCADE;
		DROP SCHEMA IF EXISTS "forum" CASCADE;

		DROP INDEX IF EXISTS demo_gin_index_ts;

		DROP COLLATION IF EXISTS case_insensitive;

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
		"name" VARCHAR(255) NOT NULL
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

		CREATE SCHEMA IF NOT EXISTS demo;

		CREATE TABLE demo.demos (
		"id"  INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"title" TEXT NOT NULL,
		"description" TEXT,
		"tags" TEXT[],
		"link" VARCHAR(255) NOT NULL,
		"user_id" UUID NOT NULL,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"upvotes" INTEGER NOT NULL DEFAULT 0,
		"downvotes" INTEGER NOT NULL DEFAULT 0,
		"thread_id" INTEGER NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE,
		"views" INTEGER NOT NULL DEFAULT 0
		);

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

		ALTER TABLE demo.demos ADD COLUMN demo_ts tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector_multilang("title"), 'A') ||
		setweight(to_tsvector_multilang(COALESCE("description", '')), 'B')
		) STORED;
		CREATE INDEX demo_gin_index_ts ON demo.demos USING GIN (demo_ts);
		`)
	if err != nil {
		panic("Error resetting demo schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO forum.topics (name) VALUES ($1) RETURNING (id)`,
		"demo").Scan(&topicID)
	if err != nil {
		panic("Error resetting demo schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".roles (name) VALUES ($1) RETURNING (id)`,
		"admin").Scan(&roleID)
	if err != nil {
		panic("Error resetting demo schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".users
		(username, display_name, email, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING (id)`,
		"mike-pech", "Mike", "test@email.com", "TestPassword", roleID).Scan(&userID)
	if err != nil {
		panic("Error resetting demo schema" + err.Error())
	}
	demo.UserID = &userID
	forumRepository = psqlRepository.NewPsqlForumRepository(testDBClient)
	demoRepository = psqlRepository.NewPsqlDemoRepository(testDBClient)
	threadSyncer = NewThreadSyncer(
		forumRepository,
		demoRepository,
		topicID,
	)
}

func TestPostThread(t *testing.T) {
	var err error

	demo.ThreadID, err = threadSyncer.PostThread(demo)
	assert.NoError(t, err)

	thread, err := forumRepository.FindThreadByID(*demo.ThreadID)
	assert.NoError(t, err)

	demoCreated, err := demoRepository.CreateDemo(demo)
	if assert.NoError(t, err) {
		assert.Equal(t, demoCreated.ThreadID, thread.ID)
		assert.Equal(t, demoCreated.Title, thread.Title)
		assert.Equal(t, demoCreated.UserID, thread.UserID)
		assert.Equal(t, demoCreated.Upvotes, thread.Upvotes)
		assert.Equal(t, demoCreated.Downvotes, thread.Downvotes)
		assert.Equal(t, demoCreated.Tags, thread.Tags)
	}
	demoUpdated.ID = demoCreated.ID
}

func TestPatchThread(t *testing.T) {
	var err error

	err = threadSyncer.PatchThread(*demoUpdated.ID, demoUpdated)
	assert.NoError(t, err)

	demoCreated, err := demoRepository.UpdateDemo(*demoUpdated.ID, demoUpdated)
	assert.NoError(t, err)

	thread, err := forumRepository.FindThreadByID(*demoCreated.ThreadID)
	if assert.NoError(t, err) {
		assert.Equal(t, demoCreated.ThreadID, thread.ID)
		assert.Equal(t, demoUpdated.Title, thread.Title)
		assert.Equal(t, demoCreated.UserID, thread.UserID)
		assert.NotEqual(t, demoCreated.UpdatedAt, thread.UpdatedAt)
		assert.Equal(t, demoCreated.Upvotes, thread.Upvotes)
		assert.Equal(t, demoCreated.Downvotes, thread.Downvotes)
		assert.Equal(t, demoCreated.Tags, thread.Tags)
	}
	if assert.NoError(t, err) {
		teardownSyncer(demoRepository, forumRepository)
	}
}

func teardownSyncer(rd *psqlRepository.PsqlDemoRepository, rf *psqlRepository.PsqlForumRepository) {
	remainderDemos, err := rd.FindDemos(nil, 0)
	if err != nil {
		panic(err)
	}
	for _, d := range *remainderDemos {
		err = rd.DeleteDemo(*d.ID)
		if err != nil {
			panic(err)
		}
	}
	err = rf.DeleteTopic(topicID)
	if err != nil {
		panic(err)
	}
}
