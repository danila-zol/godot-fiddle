package psqlRepository

import (
	"context"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"os"
	"testing"

	// "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// testDBClient     *psqlDatabase.PsqlDatabaseClient

	demoID int = 1
	// topicID          int
	// threadID         int
	// roleID           uuid.UUID
	// userID           uuid.UUID
	demoTitle        string      = "Test Demo"
	demoTitleUpdated string      = "Test UPDATE Demo"
	demoDescription  string      = "An demo for integration testing for PSQL Repo"
	demoLink         string      = "https://example.com"
	demoTags         []string    = []string{"TEST", "test"}
	demo             models.Demo = models.Demo{Title: &demoTitle, Description: &demoDescription, Link: &demoLink, ThreadID: &threadID, Tags: &demoTags}
	demoUpdated      models.Demo = models.Demo{Title: &demoTitleUpdated}
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
		os.Getenv("PSQL_CONNSTRING"), databaseConfig,
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
		"downvotes" INTEGER NOT NULL DEFAULT 0
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
		"downvotes" INTEGER NOT NULL DEFAULT 0
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
		"thread_id" INTEGER NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE
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
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO forum.threads
		(title, user_id, topic_id)
		VALUES
		($1, $2, $3)
		RETURNING
		(id)`,
		"TestDemo", userID, topicID,
	).Scan(&threadID)
	demo.ThreadID = &threadID
}

func TestCreateDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	_, err := r.CreateDemo(demo)
	assert.NoError(t, err)
}

func TestFindDemoByID(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	_, err := r.FindDemoByID(demoID)
	assert.NoError(t, err)
}

func TestFindDemoByIDNoRows(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	_, err := r.FindDemoByID(9000)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindDemosByQuery(t *testing.T) {
	var (
		demoTitleAlt       string      = "The Magnificent Seven"
		demoDescriptionAlt string      = "Marx was skint but he had sense, Engels lent him the necessary pence"
		demoTagsAlt        []string    = []string{"Cheeseboiger", "Rock the Casbah"}
		demoAlt            models.Demo = models.Demo{Title: &demoTitleAlt, Description: &demoDescriptionAlt, Link: &demoLink, ThreadID: &threadID, Tags: &demoTagsAlt, UserID: &userID}

		demoTitleAltRu       string      = "Стук"
		demoDescriptionAltRu string      = `Я скажу одно лишь слово: "Cheeseboiger"`
		demoAltRu            models.Demo = models.Demo{Title: &demoTitleAltRu, Description: &demoDescriptionAltRu, Link: &demoLink, ThreadID: &threadID, Tags: &demoTagsAlt, UserID: &userID}
	)

	r := PsqlDemoRepository{databaseClient: testDBClient}
	demoAlt.UserID = &userID
	demoAlt.ThreadID = &threadID

	for q, d := range map[string]models.Demo{"seven": demoAlt, "стук": demoAltRu} {
		resultDemo, err := r.CreateDemo(d)
		assert.NoError(t, err)

		queryDemos, err := r.FindDemosByQuery(&[]string{q})
		if assert.NoError(t, err) {
			queriedDemo := *queryDemos
			assert.Equal(t, resultDemo.Title, queriedDemo[0].Title)
			assert.Equal(t, resultDemo.Description, queriedDemo[0].Description)
			assert.Equal(t, resultDemo.Tags, queriedDemo[0].Tags)
		}
	}

	// Try to query both
	demos, err := r.FindDemosByQuery(&[]string{"cheeseboiger"})
	if assert.NoError(t, err) {
		assert.Len(t, *demos, 2)
	}
}

func TestFindDemos(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	_, err := r.FindDemos()
	assert.NoError(t, err)
}

func TestUpdateDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	resultDemo, err := r.UpdateDemo(demoID, demoUpdated)
	assert.NoError(t, err)

	modifiedDemo := demo
	modifiedDemo.ID = &demoID
	modifiedDemo.Title = &demoTitleUpdated
	modifiedDemo.CreatedAt = resultDemo.CreatedAt // Timestamps are created on DB
	modifiedDemo.UpdatedAt = resultDemo.UpdatedAt
	modifiedDemo.Upvotes = resultDemo.Upvotes
	modifiedDemo.Downvotes = resultDemo.Downvotes

	assert.Equal(t, modifiedDemo, *resultDemo)
}

func TestDeleteDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient}
	err := r.DeleteDemo(demoID)
	if assert.NoError(t, err) {
		teardownDemo(&r)
	}
}

func teardownDemo(r *PsqlDemoRepository) {
	remainderDemo, err := r.FindDemos()
	if err != nil {
		panic(err)
	}
	for _, d := range *remainderDemo {
		err = r.DeleteDemo(*d.ID)
		if err != nil {
			panic(err)
		}
	}
}
