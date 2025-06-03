package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"gamehangar/internal/domain/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"slices"
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

	// mockFileUploader mockObjectUploader
	// mockURI          string = "https://example.com"
	// mockFileInfo     os.FileInfo
	// mockFileContents []byte

	// genericUUID uuid.UUID = uuid.New()

	// notFoundResponse = `{"code":404,"message":"Not Found!"}` + "\n"

	// queryTags         = `cheeseboiger`
	// queryLimit uint64 = 1
	// queryOrder        = `newest-updated`

	demoJSON                   = `{"title":"Cool demo","description":"A very nice demo to use in your game!","userID":"` + genericUUID.String() + `"}`
	demoJSONExpected           = `{"id":1,"title":"Cool demo","description":"A very nice demo to use in your game!","userID":"` + genericUUID.String() + `","threadID":1,"key":"` + mockURI + `","thumbnailKey":"` + mockURI + `"}` + "\n"
	demoJSONExpectedMany       = `[{"id":1,"title":"Cool demo","description":"A very nice demo to use in your game!","userID":"` + genericUUID.String() + `","threadID":1,"key":"` + mockURI + `","thumbnailKey":"` + mockURI + `"}]` + "\n"
	demoJSONQueryExpected      = `[{"id":1,"title":"cheeseboiger","tags":null,"userID":"` + genericUUID.String() + `","key":"` + mockURI + `","thumbnailKey":"` + mockURI + `"},{"id":2,"title":"demo two","tags":["cheeseboiger"],"userID":"` + genericUUID.String() + `","key":null,"thumbnailKey":null}]` + "\n"
	demoJSONQueryExpectedLimit = `[{"id":1,"title":"cheeseboiger","tags":null,"userID":"` + genericUUID.String() + `","key":"` + mockURI + `","thumbnailKey":"` + mockURI + `"}]` + "\n"
	demoJSONUpdate             = `{"title":"Updated cool demo","threadID":1}`
	demoJSONUpdateExpected     = `{"id":1,"title":"Updated cool demo","description":"A very nice demo to use in your game!","userID":"` + genericUUID.String() + `","threadID":1,"key":"` + mockURI + `","thumbnailKey":"` + mockURI + `"}` + "\n"
)

func (r *mockDemoRepo) CreateDemo(demo models.Demo, demoFile, demoThumbnail io.Reader) (*models.Demo, error) {
	id := 1
	demo.ID = &id
	demo.Key = &mockURI
	demo.ThumbnailKey = &mockURI
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
func (r *mockDemoRepo) FindDemos(query []string, limit uint64, order string) (*[]models.Demo, error) {
	var (
		demoIDs    []int         = []int{1, 2, 3}
		demoTitles []string      = []string{"cheeseboiger", "demo two", "demo three"}
		demoTags   [][]string    = [][]string{nil, {"cheeseboiger"}, nil}
		demos      []models.Demo = []models.Demo{
			{ID: &demoIDs[0], Title: &demoTitles[0], Tags: &demoTags[0], UserID: &genericUUID, Key: &mockURI, ThumbnailKey: &mockURI},
			{ID: &demoIDs[1], Title: &demoTitles[1], Tags: &demoTags[1], UserID: &genericUUID},
			{ID: &demoIDs[2], Title: &demoTitles[2], Tags: &demoTags[2], UserID: &genericUUID},
		}
		resultDemos []models.Demo
	)

	if len(query) != 0 {
		for _, d := range demos {
			if *d.Title == query[0] {
				resultDemos = append(resultDemos, d)
			}
			if slices.Contains(*d.Tags, query[0]) {
				resultDemos = append(resultDemos, d)
			}
		}
	} else {
		for _, v := range r.data {
			resultDemos = append(resultDemos, v)
		}
	}
	if limit != 0 {
		resultDemos = resultDemos[:limit]
	}
	return &resultDemos, nil
}
func (r *mockDemoRepo) UpdateDemo(id int, demo models.Demo, demoFile, demoThumbnail io.Reader) (*models.Demo, error) {
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
	bodyBuffer := new(bytes.Buffer)
	mw := multipart.NewWriter(bodyBuffer) // see https://pkg.go.dev/mime/multipart
	mw.WriteField("Title", "Cool demo")
	mw.WriteField("Description", "A very nice demo to use in your game!")
	mw.WriteField("UserID", genericUUID.String())
	projPart, err := mw.CreateFormFile("demoFile", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	projPart.Write(mockFileContents)
	thumbPart, err := mw.CreateFormFile("demoThumbnail", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	thumbPart.Write(mockFileContents)
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/demos", bodyBuffer)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userTier", "freetier") // Required for attachment size check
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt, objectUploader: &mockFileUploader}

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
	bodyBuffer := new(bytes.Buffer)
	mw := multipart.NewWriter(bodyBuffer) // see https://pkg.go.dev/mime/multipart
	mw.WriteField("Title", "Updated cool demo")
	mw.WriteField("ThreadID", "1")
	projPart, err := mw.CreateFormFile("demoFile", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	projPart.Write(mockFileContents)
	thumbPart, err := mw.CreateFormFile("demoThumbnail", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	thumbPart.Write(mockFileContents)
	mw.Close()

	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/demos", bodyBuffer)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("userTier", "freetier") // Required for attachment size check
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt, objectUploader: &mockFileUploader}

	// Assertions
	if assert.NoError(t, h.PatchDemo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONUpdateExpected, rec.Body.String())
	}
}

func TestGetDemosQuery(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/demos?q=%v", queryTags), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemos(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONQueryExpected, rec.Body.String())
	}
}

func TestGetDemosQueryLimit(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/demos?q=%v&l=%v", queryTags, queryLimit), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt}

	// Assertions
	if assert.NoError(t, h.GetDemos(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, demoJSONQueryExpectedLimit, rec.Body.String())
	}
}

func TestPatchDemoNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	bodyBuffer := new(bytes.Buffer)
	mw := multipart.NewWriter(bodyBuffer) // see https://pkg.go.dev/mime/multipart
	mw.WriteField("Title", "Updated cool demo")
	mw.WriteField("ThreadID", "1")
	projPart, err := mw.CreateFormFile("demoFile", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	projPart.Write(mockFileContents)
	thumbPart, err := mw.CreateFormFile("demoThumbnail", mockFileInfo.Name())
	if err != nil {
		panic(err)
	}
	thumbPart.Write(mockFileContents)
	mw.Close()

	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/demos", bodyBuffer)
	req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("4")
	c.Set("userTier", "freetier") // Required for attachment size check
	h := &DemoHandler{logger: e.Logger, validator: v, repository: &md, syncer: &mt, objectUploader: &mockFileUploader}

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
