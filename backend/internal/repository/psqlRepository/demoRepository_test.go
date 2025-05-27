package psqlRepository

import (
	// "context"
	// "gamehangar/internal/config/psqlDatabseConfig"
	// "gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	// "gamehangar/pkg/ternMigrate"
	// "os"
	"testing"
	"time"

	// "github.com/google/uuid"
	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// independent bool = false
	// testDBClient     *psqlDatabase.PsqlDatabaseClient
	// testEnforcer *psqlCasbinClient.CasbinClient
	// testS3Client *MockS3

	demoID int = 1
	// topicID          int
	// threadID         int
	// roleID           uuid.UUID
	// userID           uuid.UUID
	demoTitle        string      = "Test Demo"
	demoTitleUpdated string      = "Test UPDATE Demo"
	demoDescription  string      = "An demo for integration testing for PSQL Repo"
	demoTags         []string    = []string{"TEST", "test"}
	demo             models.Demo = models.Demo{Title: &demoTitle, Description: &demoDescription, ThreadID: &threadID, Tags: &demoTags}
	demoUpdated      models.Demo = models.Demo{Title: &demoTitleUpdated}
)

func init() {
	if independent {
		ResetDB()
	}
	demo.UserID = &userID
	demo.ThreadID = &threadID
}

func TestCreateDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.CreateDemo(demo, nil, nil)
	assert.NoError(t, err)
}

func TestFindDemoByID(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	demo, err := r.FindDemoByID(demoID)
	if assert.NoError(t, err) { // Test view incrementation
		assert.Equal(t, uint(1), *demo.Views)
	}
}

func TestFindDemoByIDNoRows(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
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
		demoAlt            models.Demo = models.Demo{Title: &demoTitleAlt, Description: &demoDescriptionAlt, ThreadID: &threadID, Tags: &demoTagsAlt, UserID: &userID}

		demoTitleAltRu       string      = "Стук"
		demoDescriptionAltRu string      = `Я скажу одно лишь слово: "Cheeseboiger"`
		demoAltRu            models.Demo = models.Demo{Title: &demoTitleAltRu, Description: &demoDescriptionAltRu, ThreadID: &threadID, Tags: &demoTagsAlt, UserID: &userID}
	)

	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	demoAlt.UserID = &userID
	demoAlt.ThreadID = &threadID

	for q, d := range map[string]models.Demo{"seven": demoAlt, "стук": demoAltRu} {
		resultDemo, err := r.CreateDemo(d, nil, nil)
		assert.NoError(t, err)

		queryDemos, err := r.FindDemos([]string{q}, 0, "highest-rated")
		if assert.NoError(t, err) {
			queriedDemo := *queryDemos
			assert.Equal(t, resultDemo.Title, queriedDemo[0].Title)
			assert.Equal(t, resultDemo.Description, queriedDemo[0].Description)
			assert.Equal(t, resultDemo.Tags, queriedDemo[0].Tags)
		}
	}

	// Try to query both and check ordering
	demos, err := r.FindDemos([]string{"cheeseboiger"}, 0, "newest-updated")
	if assert.NoError(t, err) {
		d := *demos
		assert.Len(t, d, 2)
		var timeOrder, timeOrderExpected []time.Time
		timeOrderExpected = []time.Time{*d[0].UpdatedAt, *d[1].UpdatedAt}
		for _, m := range d {
			timeOrder = append(timeOrder, *m.UpdatedAt)
		}
		assert.Equal(
			t,
			timeOrderExpected,
			timeOrder,
		)
	}
	// Query with limit
	demos, err = r.FindDemos([]string{"cheeseboiger"}, 1, "newest-updated")
	if assert.NoError(t, err) {
		assert.Len(t, *demos, 1)
	}
}

func TestFindDemos(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindDemos(nil, 0, "")
	assert.NoError(t, err)
}

func TestUpdateDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}

	oldDemo, err := r.FindDemoByID(demoID)
	assert.NoError(t, err)

	resultDemo, err := r.UpdateDemo(demoID, demoUpdated, nil, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, oldDemo.CreatedAt, resultDemo.CreatedAt)
		assert.Equal(t, oldDemo.Rating, resultDemo.Rating)
		assert.Equal(t, oldDemo.Views, resultDemo.Views)

		assert.NotEqual(t, oldDemo.UpdatedAt, resultDemo.UpdatedAt)

		assert.NotEqual(t, oldDemo.Title, resultDemo.Title)
		assert.Equal(t, demoUpdated.Title, resultDemo.Title)
	}
}

func TestDeleteDemo(t *testing.T) {
	r := PsqlDemoRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	err := r.DeleteDemo(demoID)
	if assert.NoError(t, err) {
		teardownDemo(&r)
	}
}

func teardownDemo(r *PsqlDemoRepository) {
	remainderDemo, err := r.FindDemos(nil, 0, "")
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
