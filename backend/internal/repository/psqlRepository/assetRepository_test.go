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
	independent  bool = false
	testDBClient *psqlDatabase.PsqlDatabaseClient

	views uint = 0

	assetID          int          = 1
	assetName        string       = "Test Asset"
	assetNameUpdated string       = "Test UPDATE Asset"
	assetDescription string       = "An asset for integration testing for PSQL Repo"
	assetLink        string       = "https://example.com"
	assetVersion     int          = 1
	asset            models.Asset = models.Asset{Name: &assetName, Description: &assetDescription, Link: &assetLink}
	assetUpdated     models.Asset = models.Asset{Name: &assetNameUpdated, Version: &assetVersion}
)

func ResetDB() {
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
		DROP TRIGGER IF EXISTS increment_asset_version_on_update ON asset.assets; 
		DROP TRIGGER IF EXISTS increment_topic_version_on_update ON forum.topics; 
		DROP TRIGGER IF EXISTS increment_role_version_on_update ON "user".roles; 

		DROP COLLATION IF EXISTS case_insensitive;

		DROP INDEX IF EXISTS demo_gin_index_ts;
		DROP INDEX IF EXISTS asset_gin_index_ts;
		DROP INDEX IF EXISTS thread_gin_index_ts;
		DROP INDEX IF EXISTS message_gin_index_ts;

		DROP SCHEMA IF EXISTS "asset" CASCADE;
		DROP SCHEMA IF EXISTS "demo" CASCADE;
		DROP SCHEMA IF EXISTS "user" CASCADE;
		DROP SCHEMA IF EXISTS "forum" CASCADE;

		CREATE SCHEMA IF NOT EXISTS demo;
		CREATE SCHEMA IF NOT EXISTS forum;
		CREATE SCHEMA IF NOT EXISTS "user";
		CREATE SCHEMA IF NOT EXISTS asset;

		CREATE TABLE "user".roles (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"name" VARCHAR(255) NOT NULL,
		-- "permissions" VARCHAR(64)[]
		"version" INTEGER NOT NULL DEFAULT 1
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
		"upvotes" INTEGER NOT NULL DEFAULT 1,
		"downvotes" INTEGER NOT NULL DEFAULT 1,
		"rating" DECIMAL GENERATED ALWAYS AS (upvotes::DECIMAL / downvotes) STORED,
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
		"upvotes" INTEGER NOT NULL DEFAULT 1,
		"downvotes" INTEGER NOT NULL DEFAULT 1,
		"rating" DECIMAL GENERATED ALWAYS AS (upvotes::DECIMAL / downvotes) STORED,
		"views" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE asset.assets (
		"id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"name" TEXT NOT NULL,
		"description" TEXT,
		"link" VARCHAR(255) NOT NULL,
		"tags" TEXT[],
		"version" INTEGER NOT NULL DEFAULT 1,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"upvotes" INTEGER NOT NULL DEFAULT 1,
		"downvotes" INTEGER NOT NULL DEFAULT 1,		
		"rating" DECIMAL GENERATED ALWAYS AS (upvotes::DECIMAL / downvotes) STORED,
		"views" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE demo.demos (
		"id"  INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"title" TEXT NOT NULL,
		"description" TEXT,
		"tags" TEXT[],
		"link" VARCHAR(255) NOT NULL,
		"user_id" UUID NOT NULL,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"upvotes" INTEGER NOT NULL DEFAULT 1,
		"downvotes" INTEGER NOT NULL DEFAULT 1,
		"rating" DECIMAL GENERATED ALWAYS AS (upvotes::DECIMAL / downvotes) STORED,
		"views" INTEGER NOT NULL DEFAULT 0,
		"thread_id" INTEGER NOT NULL REFERENCES forum.threads (id) ON DELETE CASCADE
		);

		CREATE OR REPLACE FUNCTION increment_version()
		RETURNS TRIGGER AS
		$func$
		BEGIN
		NEW.version := OLD.version + 1;
		RETURN NEW;
		END;
		$func$ LANGUAGE plpgsql;

		CREATE TRIGGER increment_asset_version_on_update
		BEFORE UPDATE ON asset.assets
		FOR EACH ROW
		WHEN ((OLD.name IS DISTINCT FROM NEW.name) 
		OR (OLD.description IS DISTINCT FROM NEW.description)
		OR (OLD.link IS DISTINCT FROM NEW.link)
		OR (OLD.tags IS DISTINCT FROM NEW.tags))
		EXECUTE FUNCTION increment_version();

		CREATE TRIGGER increment_topic_version_on_update
		BEFORE UPDATE ON forum.topics
		FOR EACH ROW
		EXECUTE FUNCTION increment_version();

		CREATE TRIGGER increment_role_version_on_update
		BEFORE UPDATE ON "user".roles
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

		ALTER TABLE asset.assets ADD COLUMN asset_ts tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector_multilang("name"), 'A') ||
		setweight(to_tsvector_multilang(COALESCE("description", '')), 'B')
		) STORED;
		CREATE INDEX asset_gin_index_ts ON asset.assets USING GIN (asset_ts);

		ALTER TABLE demo.demos ADD COLUMN demo_ts tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector_multilang("title"), 'A') ||
		setweight(to_tsvector_multilang(COALESCE("description", '')), 'B')
		) STORED;
		CREATE INDEX demo_gin_index_ts ON demo.demos USING GIN (demo_ts);

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
		panic("Error resetting assets schema" + err.Error())
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
}

func init() {
	if independent || 1 == 1 { // HACK!
		ResetDB()
	}
}

func TestCreateAsset(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.CreateAsset(asset)
	assert.NoError(t, err)
}

func TestFindAssetByID(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	asset, err := r.FindAssetByID(assetID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *asset.Views)
	}
}

func TestFindAssetByIDNoRows(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindAssetByID(39)
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestFindAssets(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindAssets(nil, 0, "")
	assert.NoError(t, err)

}

func TestFindAssetsByQuery(t *testing.T) {
	var (
		assetTitleAlt       string       = "The Magnificent Seven"
		assetDescriptionAlt string       = `Marx was skint but he had sense, Engels lent him the necessary pence`
		assetTagsAlt        []string     = []string{"Cheeseboiger", "Rock the Casbah"}
		assetAlt            models.Asset = models.Asset{Name: &assetTitleAlt, Description: &assetDescriptionAlt, Link: &assetLink, Tags: &assetTagsAlt}

		assetTitleAltRu       string       = "Стук"
		assetDescriptionAltRu string       = `Я скажу одно лишь слово: "Cheeseboiger"`
		assetAltRu            models.Asset = models.Asset{Name: &assetTitleAltRu, Description: &assetDescriptionAltRu, Link: &assetLink}
	)

	r := PsqlAssetRepository{databaseClient: testDBClient}

	for q, d := range map[string]models.Asset{"seven": assetAlt, "стук": assetAltRu} {
		resultAsset, err := r.CreateAsset(d)
		assert.NoError(t, err)

		queryAssets, err := r.FindAssets([]string{q}, 0, "")
		if assert.NoError(t, err) {
			queriedAsset := *queryAssets
			assert.Equal(t, resultAsset.Name, queriedAsset[0].Name)
			assert.Equal(t, resultAsset.Description, queriedAsset[0].Description)
		}
	}

	// Try to query both and check ordering
	assets, err := r.FindAssets([]string{"cheeseboiger"}, 0, "newest-updated")
	if assert.NoError(t, err) {
		a := *assets
		assert.Len(t, a, 2)
		var timeOrder, timeOrderExpected []time.Time
		timeOrderExpected = []time.Time{*a[0].UpdatedAt, *a[1].UpdatedAt}
		for _, m := range a {
			timeOrder = append(timeOrder, *m.UpdatedAt)
		}
		assert.Equal(
			t,
			timeOrderExpected,
			timeOrder,
		)
	}
	// Query with limit
	assets, err = r.FindAssets([]string{"cheeseboiger"}, 1, "")
	if assert.NoError(t, err) {
		assert.Len(t, *assets, 1)
	}
}

func TestUpdateAsset(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}

	oldAsset, err := r.FindAssetByID(assetID)
	assert.NoError(t, err)

	resultAsset, err := r.UpdateAsset(assetID, assetUpdated)
	if assert.NoError(t, err) {
		assert.Equal(t, oldAsset.CreatedAt, resultAsset.CreatedAt)
		assert.Equal(t, oldAsset.Link, resultAsset.Link)
		assert.Equal(t, oldAsset.Rating, resultAsset.Rating)
		assert.Equal(t, oldAsset.Views, resultAsset.Views)

		assert.NotEqual(t, oldAsset.UpdatedAt, resultAsset.UpdatedAt)

		assert.NotEqual(t, oldAsset.Name, resultAsset.Name)
		assert.Equal(t, assetUpdated.Name, resultAsset.Name)

		newVersion := *assetUpdated.Version + 1
		assert.Equal(t, newVersion, *resultAsset.Version)
	}
}

func TestUpdateAssetMultiple(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}

	for i := 2; i < 6; i++ {
		newName := *assetUpdated.Name + " New"
		assetUpdated.Name = &newName
		newVersion := i
		assetUpdated.Version = &newVersion

		oldAsset, err := r.FindAssetByID(assetID)
		assert.NoError(t, err)

		resultAsset, err := r.UpdateAsset(assetID, assetUpdated)
		if assert.NoError(t, err) {
			assert.Equal(t, oldAsset.CreatedAt, resultAsset.CreatedAt)
			assert.Equal(t, oldAsset.Link, resultAsset.Link)
			assert.Equal(t, oldAsset.Rating, resultAsset.Rating)
			assert.Equal(t, oldAsset.Views, resultAsset.Views)

			assert.NotEqual(t, oldAsset.UpdatedAt, resultAsset.UpdatedAt)

			assert.NotEqual(t, oldAsset.Name, resultAsset.Name)
			assert.Equal(t, assetUpdated.Name, resultAsset.Name)

			newVersion := *assetUpdated.Version + 1
			assert.Equal(t, newVersion, *resultAsset.Version)
		}
	}
}

func TestUpdateAssetConflict(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.UpdateAsset(assetID, assetUpdated)
	if assert.Error(t, err) {
		assert.Equal(t, r.conflictErr, err)
	}
}

func TestDeleteAsset(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	err := r.DeleteAsset(assetID)
	if assert.NoError(t, err) {
		teardownAsset(&r)
	}
}

func teardownAsset(r *PsqlAssetRepository) {
	remainderAssets, err := r.FindAssets(nil, 0, "")
	if err != nil {
		panic(err)
	}
	for _, a := range *remainderAssets {
		err = r.DeleteAsset(*a.ID)
		if err != nil {
			panic(err)
		}
	}
}
