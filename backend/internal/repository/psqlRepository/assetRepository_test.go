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

func init() {
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../../.env") // Whatever...
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
		DROP SCHEMA IF EXISTS "asset" CASCADE;

		DROP COLLATION IF EXISTS case_insensitive;
		DROP INDEX IF EXISTS asset_gin_index_ts;

		CREATE SCHEMA IF NOT EXISTS asset;

		CREATE TABLE asset.assets (
		"id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"name" TEXT NOT NULL,
		"description" TEXT,
		"link" VARCHAR(255) NOT NULL,
		"tags" TEXT[],
		"version" INTEGER NOT NULL DEFAULT 1,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
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

		CREATE TRIGGER increment_asset_version_on_update
		BEFORE UPDATE ON asset.assets
		FOR EACH ROW
		WHEN ((OLD.name IS DISTINCT FROM NEW.name) 
		OR (OLD.description IS DISTINCT FROM NEW.description)
		OR (OLD.link IS DISTINCT FROM NEW.link)
		OR (OLD.tags IS DISTINCT FROM NEW.tags))
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
		`)
	if err != nil {
		panic("Error resetting assets schema" + err.Error())
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
	_, err := r.FindAssets(nil, 0)
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

		queryAssets, err := r.FindAssets([]string{q}, 0)
		if assert.NoError(t, err) {
			queriedAsset := *queryAssets
			assert.Equal(t, resultAsset.Name, queriedAsset[0].Name)
			assert.Equal(t, resultAsset.Description, queriedAsset[0].Description)
		}
	}

	// Try to query both and check ordering
	assets, err := r.FindAssets([]string{"cheeseboiger"}, 0)
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
	assets, err = r.FindAssets([]string{"cheeseboiger"}, 1)
	if assert.NoError(t, err) {
		assert.Len(t, *assets, 1)
	}
}

func TestUpdateAsset(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	resultAsset, err := r.UpdateAsset(assetID, assetUpdated)
	assert.NoError(t, err)

	modifiedAsset := asset // Manual update
	modifiedAsset.ID = &assetID
	modifiedAsset.Name = &assetNameUpdated
	modifiedAsset.CreatedAt = resultAsset.CreatedAt // Timestamps are created on DB
	modifiedAsset.UpdatedAt = resultAsset.UpdatedAt
	modifiedAsset.Views = resultAsset.Views
	newVersion := *assetUpdated.Version + 1
	modifiedAsset.Version = &newVersion

	assert.Equal(t, modifiedAsset, *resultAsset)
}

func TestUpdateAssetMultiple(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}

	modifiedAsset := asset // Manual update
	modifiedAsset.ID = &assetID

	for i := 2; i < 6; i++ {
		newName := *assetUpdated.Name + " New"
		assetUpdated.Name = &newName
		newVersion := i
		assetUpdated.Version = &newVersion

		resultAsset, err := r.UpdateAsset(assetID, assetUpdated)
		assert.NoError(t, err)

		newerVersion := i + 1
		modifiedAsset.Version = &newerVersion
		modifiedAsset.Name = assetUpdated.Name
		modifiedAsset.CreatedAt = resultAsset.CreatedAt // Timestamps are created on DB
		modifiedAsset.UpdatedAt = resultAsset.UpdatedAt
		modifiedAsset.Views = resultAsset.Views

		assert.Equal(t, modifiedAsset, *resultAsset)
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
	remainderAssets, err := r.FindAssets(nil, 0)
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
