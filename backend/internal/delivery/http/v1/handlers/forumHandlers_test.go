package handlers

import (
	"errors"
	"fmt"
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

type mockForumRepo struct {
	topicData   map[int]models.Topic
	threadData  map[int]models.Thread
	messageData map[int]models.Message
	notFoundErr error
	conflictErr error
}

var (
	// v = validator.New(validator.WithRequiredStructEnabled())
	mf = mockForumRepo{
		topicData:   make(map[int]models.Topic, 1),
		threadData:  make(map[int]models.Thread, 1),
		messageData: make(map[int]models.Message, 1),
		notFoundErr: errors.New("Not Found"),
		conflictErr: errors.New("Record conflict!"),
	}

	// genericUUID uuid.UUID = uuid.New()

	// notFoundResponse = `{"code":404,"message":"Not Found!"}` + "\n"
	// conflictResponse = `{"code":409,"message":"Error: unable to update the record due to an edit conflict, please try again!"}` + "\n"

	// query             = `cheeseboiger`
	// queryLimit uint64 = 1

	topicJSON               = `{"name":"Cool topic"}`
	topicJSONExpected       = `{"id":1,"name":"Cool topic","version":1}` + "\n"
	topicJSONExpectedMany   = `[{"id":1,"name":"Cool topic","version":1}]` + "\n"
	topicJSONUpdate         = `{"name":"Updated cool topic","version":1}`
	topicJSONUpdateInvalid  = `{"name":"Updated cool topic"}`
	topicJSONUpdateExpected = `{"id":1,"name":"Updated cool topic","version":2}` + "\n"

	threadJSON                   = `{"title":"Cool Thread","userID":"` + genericUUID.String() + `","topicID":1}`
	threadJSONExpected           = `{"id":1,"title":"Cool Thread","userID":"` + genericUUID.String() + `","topicID":1}` + "\n"
	threadJSONExpectedMany       = `[{"id":1,"title":"Cool Thread","userID":"` + genericUUID.String() + `","topicID":1}]` + "\n"
	threadJSONQueryExpected      = `[{"id":1,"title":"cheeseboiger","userID":"` + genericUUID.String() + `","topicID":2,"tags":null},{"id":2,"title":"thread two","userID":"` + genericUUID.String() + `","topicID":2,"tags":["cheeseboiger"]}]` + "\n"
	threadJSONQueryExpectedLimit = `[{"id":1,"title":"cheeseboiger","userID":"` + genericUUID.String() + `","topicID":2,"tags":null}]` + "\n"
	threadJSONUpdate             = `{"title":"Updated cool Thread"}`
	threadJSONUpdateExpected     = `{"id":1,"title":"Updated cool Thread","userID":"` + genericUUID.String() + `","topicID":1}` + "\n"

	messageJSON                   = `{"title":"Cool message","userID":"` + genericUUID.String() + `","threadID":1}`
	messageJSONExpected           = `{"id":1,"threadID":1,"userID":"` + genericUUID.String() + `","title":"Cool message"}` + "\n"
	messageJSONExpectedMany       = `[{"id":1,"threadID":1,"userID":"` + genericUUID.String() + `","title":"Cool message"}]` + "\n"
	messageJSONQueryExpected      = `[{"id":1,"threadID":2,"userID":"` + genericUUID.String() + `","title":"cheeseboiger","tags":null},{"id":2,"threadID":2,"userID":"` + genericUUID.String() + `","title":"message two","tags":["cheeseboiger"]}]` + "\n"
	messageJSONQueryExpectedLimit = `[{"id":1,"threadID":2,"userID":"` + genericUUID.String() + `","title":"cheeseboiger","tags":null}]` + "\n"
	messageJSONUpdate             = `{"title":"Updated cool message"}`
	messageJSONUpdateExpected     = `{"id":1,"threadID":1,"userID":"` + genericUUID.String() + `","title":"Updated cool message"}` + "\n"
)

func (r *mockForumRepo) CreateTopic(topic models.Topic) (*models.Topic, error) {
	id := 1
	topic.ID = &id
	version := 1
	topic.Version = &version
	r.topicData[id] = topic
	resultTopic, ok := r.topicData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultTopic, nil
}
func (r *mockForumRepo) FindTopicByID(id int) (*models.Topic, error) {
	topic, ok := r.topicData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &topic, nil
}
func (r *mockForumRepo) FindTopics() (*[]models.Topic, error) {
	var t []models.Topic
	for _, v := range r.topicData {
		t = append(t, v)
	}
	return &t, nil
}
func (r *mockForumRepo) UpdateTopic(id int, topic models.Topic) (*models.Topic, error) {
	var resultTopic models.Topic
	_, ok := r.topicData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	resultTopic = r.topicData[id]
	if *resultTopic.Version != *topic.Version {
		return nil, r.ConflictErr()
	}
	if topic.Name != nil {
		resultTopic.Name = topic.Name
		n := *topic.Version + 1
		resultTopic.Version = &n
		r.topicData[id] = resultTopic
	}
	resultTopic = r.topicData[id]
	return &resultTopic, nil
}
func (r *mockForumRepo) DeleteTopic(id int) error {
	_, ok := r.topicData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.topicData, id)
	return nil
}

func (r *mockForumRepo) CreateThread(thread models.Thread) (*models.Thread, error) {
	id := 1
	thread.ID = &id
	r.threadData[id] = thread
	resultThread, ok := r.threadData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultThread, nil
}
func (r *mockForumRepo) FindThreads(query []string, limit uint64) (*[]models.Thread, error) {
	var (
		topicID      int             = 2
		threadIDs    []int           = []int{1, 2, 3}
		threadTitles []string        = []string{"cheeseboiger", "thread two", "thread three"}
		threadTags   [][]string      = [][]string{nil, []string{"cheeseboiger"}, nil}
		threads      []models.Thread = []models.Thread{
			{ID: &threadIDs[0], TopicID: &topicID, UserID: &genericUUID, Title: &threadTitles[0], Tags: &threadTags[0]},
			{ID: &threadIDs[1], TopicID: &topicID, UserID: &genericUUID, Title: &threadTitles[1], Tags: &threadTags[1]},
			{ID: &threadIDs[2], TopicID: &topicID, UserID: &genericUUID, Title: &threadTitles[2], Tags: &threadTags[2]},
		}
		resultThreads []models.Thread
	)
	if len(query) != 0 {
		for _, t := range threads {
			if *t.Title == query[0] {
				resultThreads = append(resultThreads, t)
			}
			if slices.Contains(*t.Tags, query[0]) {
				resultThreads = append(resultThreads, t)
			}
		}
	} else {
		for _, v := range r.threadData {
			resultThreads = append(resultThreads, v)
		}
	}
	if limit != 0 {
		resultThreads = resultThreads[:limit]
	}
	return &resultThreads, nil
}
func (r *mockForumRepo) FindThreadByID(id int) (*models.Thread, error) {
	thread, ok := r.threadData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &thread, nil
}
func (r *mockForumRepo) UpdateThread(id int, thread models.Thread) (*models.Thread, error) {
	var resultThread models.Thread
	_, ok := r.threadData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	resultThread = r.threadData[id]
	if thread.Title != nil {
		resultThread.Title = thread.Title
		r.threadData[id] = resultThread
	}
	resultThread = r.threadData[id]
	return &resultThread, nil
}
func (r *mockForumRepo) DeleteThread(id int) error {
	_, ok := r.threadData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.threadData, id)
	return nil
}

func (r *mockForumRepo) CreateMessage(message models.Message) (*models.Message, error) {
	id := 1
	message.ID = &id
	r.messageData[id] = message
	resultmessage, ok := r.messageData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultmessage, nil
}
func (r *mockForumRepo) FindMessageByID(id int) (*models.Message, error) {
	message, ok := r.messageData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &message, nil
}
func (r *mockForumRepo) FindMessages(query []string, limit uint64) (*[]models.Message, error) {
	var (
		topicID       int              = 2
		messageIDs    []int            = []int{1, 2, 3}
		messageTitles []string         = []string{"cheeseboiger", "message two", "message three"}
		messageTags   [][]string       = [][]string{nil, []string{"cheeseboiger"}, nil}
		messages      []models.Message = []models.Message{
			{ID: &messageIDs[0], ThreadID: &topicID, UserID: &genericUUID, Title: &messageTitles[0], Tags: &messageTags[0]},
			{ID: &messageIDs[1], ThreadID: &topicID, UserID: &genericUUID, Title: &messageTitles[1], Tags: &messageTags[1]},
			{ID: &messageIDs[2], ThreadID: &topicID, UserID: &genericUUID, Title: &messageTitles[2], Tags: &messageTags[2]},
		}
		resultMessages []models.Message
	)
	if len(query) != 0 {
		for _, m := range messages {
			if *m.Title == query[0] {
				resultMessages = append(resultMessages, m)
			}
			if slices.Contains(*m.Tags, query[0]) {
				resultMessages = append(resultMessages, m)
			}
		}
	} else {
		for _, m := range r.messageData {
			resultMessages = append(resultMessages, m)
		}
	}
	if limit != 0 {
		resultMessages = resultMessages[:limit]
	}
	return &resultMessages, nil
}
func (r *mockForumRepo) FindMessagesByThreadID(threadID int) (*[]models.Message, error) {
	var (
		messageIDs    []int            = []int{1, 2}
		messageTitles []string         = []string{"message one", "message two"}
		t             []models.Message = []models.Message{
			{ID: &messageIDs[0], ThreadID: &threadID, UserID: &genericUUID, Title: &messageTitles[0]},
			{ID: &messageIDs[1], ThreadID: &threadID, UserID: &genericUUID, Title: &messageTitles[1]},
		}
	)
	if threadID != 1 {
		return nil, r.NotFoundErr()
	}
	return &t, nil
}
func (r *mockForumRepo) UpdateMessage(id int, message models.Message) (*models.Message, error) {
	var resultMessage models.Message
	_, ok := r.messageData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	resultMessage = r.messageData[id]
	if message.Title != nil {
		resultMessage.Title = message.Title
		r.messageData[id] = resultMessage
	}
	resultMessage = r.messageData[id]
	return &resultMessage, nil
}
func (r *mockForumRepo) DeleteMessage(id int) error {
	_, ok := r.messageData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.messageData, id)
	return nil
}

func (r *mockForumRepo) NotFoundErr() error { return r.notFoundErr }
func (r *mockForumRepo) ConflictErr() error { return r.conflictErr }

func TestPostTopic(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/topics", strings.NewReader(topicJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PostTopic(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, topicJSONExpected, rec.Body.String())
	}
}

func TestGetTopicByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopicByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, topicJSONExpected, rec.Body.String())
	}
}

func TestGetTopicByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopicByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestGetTopics(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopics(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, topicJSONExpectedMany, rec.Body.String())
	}
}

func TestPatchTopic(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/topics", strings.NewReader(topicJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchTopic(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, topicJSONUpdateExpected, rec.Body.String())
	}
}

func TestPatchTopicNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/topics", strings.NewReader(topicJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchTopic(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPatchTopicUnprocessable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/topics", strings.NewReader(topicJSONUpdateInvalid))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchTopic(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestPatchTopicConflict(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/topics", strings.NewReader(topicJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchTopic(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, conflictResponse, rec.Body.String())
	}
}

func TestDeleteTopic(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopicByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteTopicUnprocesable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9bc3c90")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopicByID(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteTopicNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetTopicByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPostThread(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/threads", strings.NewReader(threadJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PostThread(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, threadJSONExpected, rec.Body.String())
	}
}

func TestGetThreadByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreadByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadJSONExpected, rec.Body.String())
	}
}

func TestGetThreadByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreadByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestGetThreadsQuery(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/threads?q="+query, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreads(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadJSONQueryExpected, rec.Body.String())
	}
}

func TestGetThreadsQueryLimit(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/threads?q=%v&l=%v", query, queryLimit), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreads(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadJSONQueryExpectedLimit, rec.Body.String())
	}
}

func TestGetThreads(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreads(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadJSONExpectedMany, rec.Body.String())
	}
}

func TestPatchThread(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/threads", strings.NewReader(threadJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchThread(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadJSONUpdateExpected, rec.Body.String())
	}
}

func TestPatchThreadNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/threads", strings.NewReader(threadJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchThread(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestDeleteThread(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreadByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteThreadUnprocesable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9bc3c90")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreadByID(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteThreadNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/threads", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetThreadByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPostMessage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/messages", strings.NewReader(messageJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PostMessage(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, messageJSONExpected, rec.Body.String())
	}
}

func TestGetMessageByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessageByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, messageJSONExpected, rec.Body.String())
	}
}

func TestGetMessageByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessageByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestGetMessagesQuery(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/messages?q="+query, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessages(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, messageJSONQueryExpected, rec.Body.String())
	}
}

func TestGetMessagesQueryLimit(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/messages?q=%v&l=%v", query, queryLimit), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessages(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, messageJSONQueryExpectedLimit, rec.Body.String())
	}
}

func TestGetMessages(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessages(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, messageJSONExpectedMany, rec.Body.String())
	}
}

func TestPatchMessage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/messages", strings.NewReader(messageJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchMessage(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, messageJSONUpdateExpected, rec.Body.String())
	}
}

func TestPatchMessageNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/messages", strings.NewReader(messageJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.PatchMessage(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestDeleteMessage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessageByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteMessageUnprocesable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9bc3c90")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessageByID(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteMessageNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93")
	h := &ForumHandler{logger: e.Logger, validator: v, repository: &mf}

	// Assertions
	if assert.NoError(t, h.GetMessageByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}
