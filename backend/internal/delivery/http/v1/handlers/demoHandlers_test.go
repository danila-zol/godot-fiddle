package handlers

import (
	"errors"
	"gamehangar/internal/domain/models"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	// "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockThreadSyncer struct {
	topicID int
}

type mockDemoRepo struct {
	data        map[int]models.Demo
	notFoundErr error
}

var (
	// v  = validator.New(validator.WithRequiredStructEnabled())
	mt = mockThreadSyncer{topicID: 1}
	md = mockDemoRepo{
		data:        make(map[int]models.Demo, 1),
		notFoundErr: errors.New("Not Found"),
	}

	// genericUUID string = "9c6ac0b1-b97e-4356-a6e1-dc6b52324220"

	// notFoundResponse = `{"code":404,"message":"Not Found!"}` + "\n"

	demoJSON               = `{"title":"Cool demo","description":"A very nice demo to use in your game!","link":"https://example.com","userID":"` + genericUUID + `"}`
	demoJSONExpected       = `{"id":1,"title":"Cool demo","description":"A very nice demo to use in your game!","link":"https://example.com","userID":"` + genericUUID + `","threadID":1}` + "\n"
	demoJSONExpectedMany   = `[{"id":1,"title":"Cool demo","description":"A very nice demo to use in your game!","link":"https://example.com","userID":"` + genericUUID + `","threadID":1}]` + "\n"
	demoQuery              = `cheeseboiger`
	demoJSONQueryExpected  = `[{"id":1,"title":"cheeseboiger","link":"link.com","tags":null,"userID":"` + genericUUID + `"},{"id":2,"title":"demo two","link":"example.com","tags":["cheeseboiger"],"userID":"` + genericUUID + `"}]` + "\n"
	demoJSONUpdate         = `{"title":"Updated cool demo","threadID":1}`
	demoJSONUpdateExpected = `{"id":1,"title":"Updated cool demo","description":"A very nice demo to use in your game!","link":"https://example.com","userID":"` + genericUUID + `","threadID":1}` + "\n"
)

func (r *mockDemoRepo) CreateDemo(demo models.Demo) (*models.Demo, error) {
	id := 1
	demo.ID = &id
	r.data[id] = demo
	return &demo, nil
}
func (r *mockDemoRepo) FindDemoByID(id int) (*models.Demo, error) {
	a, ok := r.data[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &a, nil
}
func (r *mockDemoRepo) FindDemos() (*[]models.Demo, error) {
	var a []models.Demo
	for _, v := range r.data {
		a = append(a, v)
	}
	return &a, nil
}
func (r *mockDemoRepo) FindDemosByQuery(query *[]string) (*[]models.Demo, error) {
	var (
		demoIDs    []int         = []int{1, 2, 3}
		demoTitles []string      = []string{"cheeseboiger", "demo two", "demo three"}
		demoLinks  []string      = []string{"link.com", "example.com", "e.com"}
		demoTags   [][]string    = [][]string{nil, []string{"cheeseboiger"}, nil}
		demos      []models.Demo = []models.Demo{
			{ID: &demoIDs[0], Title: &demoTitles[0], Link: &demoLinks[0], Tags: &demoTags[0], UserID: &genericUUID},
			{ID: &demoIDs[1], Title: &demoTitles[1], Link: &demoLinks[1], Tags: &demoTags[1], UserID: &genericUUID},
			{ID: &demoIDs[2], Title: &demoTitles[2], Link: &demoLinks[2], Tags: &demoTags[2], UserID: &genericUUID},
		}
		resultDemos []models.Demo
	)
	q := *query
	for _, d := range demos {
		if *d.Title == q[0] {
			resultDemos = append(resultDemos, d)
		}
		if slices.Contains(*d.Tags, q[0]) {
			resultDemos = append(resultDemos, d)
		}
	}
	return &resultDemos, nil
}
func (r *mockDemoRepo) UpdateDemo(id int, demo models.Demo) (*models.Demo, error) {
	var a models.Demo
	_, ok := r.data[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	a = r.data[id]
	if demo.Title != nil {
		a.Title = demo.Title
		r.data[id] = a
	}
	a = r.data[id]
	return &a, nil
}
func (r *mockDemoRepo) DeleteDemo(id int) error {
	_, ok := r.data[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.data, id)
	return nil
}
func (r *mockDemoRepo) NotFoundErr() error { return r.notFoundErr }

func (s *mockThreadSyncer) PostThread(demo models.Demo) (*int, error) {
	threadID := 1
	return &threadID, nil
}
func (s *mockThreadSyncer) PatchThread(demoID int, demo models.Demo) error { return nil }

func TestPostDemo(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/demos", strings.NewReader(demoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.PostDemo(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, demoJSONExpected, rec.Body.String())
	}
}

func TestGetDemoByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/demos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemoById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONExpected, rec.Body.String())
	}
}

func TestGetDemoByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/demos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("4")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemoById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestGetDemos(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/demos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemos(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONExpectedMany, rec.Body.String())
	}
}

func TestPatchDemo(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/demos", strings.NewReader(demoJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.PatchDemo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONUpdateExpected, rec.Body.String())
	}
}

func TestGetDemosByQuery(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/demos?q="+demoQuery, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemos(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONQueryExpected, rec.Body.String())
	}
}

func TestPatchDemoNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/demos", strings.NewReader(demoJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("4")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.PatchDemo(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestDeleteDemo(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/demos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemoById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteDemoNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/demos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("5")
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemoById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}
