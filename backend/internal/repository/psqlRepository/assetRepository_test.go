package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	testDBClient     *psqlDatabase.PsqlDatabaseClient
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
		os.Getenv("PSQL_CONNSTRING"), databaseConfig,
	)
	if err != nil {
		panic("Error setting up new DatabaseClient")
	}
	c, _ := testDBClient.AcquireConn() // WARNING! Integration tests DROP TABLEs
	_, err = c.Exec(context.Background(), `
		DROP TRIGGER IF EXISTS increment_asset_version_on_update ON asset.assets; 
		DROP SCHEMA IF EXISTS "asset" CASCADE;

		CREATE SCHEMA IF NOT EXISTS asset;

		CREATE TABLE asset.assets(
		"id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		"name" VARCHAR(255) NOT NULL,
		"description" VARCHAR,
		"link" VARCHAR(255) NOT NULL,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"version" INTEGER NOT NULL DEFAULT 1);

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
		EXECUTE FUNCTION increment_version();
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
	_, err := r.FindAssetByID(assetID)
	assert.NoError(t, err)
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
	_, err := r.FindAssets()
	assert.NoError(t, err)

}

func TestUpdateAsset(t *testing.T) {
	r := PsqlAssetRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	resultAsset, err := r.UpdateAsset(assetID, assetUpdated)
	assert.NoError(t, err)

	modifiedAsset := asset // Manual update
	modifiedAsset.ID = &assetID
	modifiedAsset.Name = &assetNameUpdated
	modifiedAsset.CreatedAt = resultAsset.CreatedAt // Timestamps are created on DB
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
	assert.NoError(t, err)
}
