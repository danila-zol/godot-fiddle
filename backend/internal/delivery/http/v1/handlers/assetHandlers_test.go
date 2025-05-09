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

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockAssetRepo struct {
	data        map[int]models.Asset
	notFoundErr error
	conflictErr error
}

var (
	v  = validator.New(validator.WithRequiredStructEnabled())
	ma = mockAssetRepo{
		data:        make(map[int]models.Asset, 1),
		notFoundErr: errors.New("Not Found"),
		conflictErr: errors.New("Record conflict!"),
	}

	// notFoundResponse = `{"code":404,"message":"Not Found!"}` + "\n"
	// conflictResponse = `{"code":409,"message":"Error: unable to update the record due to an edit conflict, please try again!"}` + "\n"

	// query                        = `cheeseboiger`
	// queryLimit                 uint64 = 1

	assetJSON                   = `{"name":"Cool asset","description":"A very nice asset to use in your game!","link":"https://example.com"}`
	assetJSONExpected           = `{"id":1,"name":"Cool asset","description":"A very nice asset to use in your game!","link":"https://example.com","version":1}` + "\n"
	assetJSONExpectedMany       = `[{"id":1,"name":"Cool asset","description":"A very nice asset to use in your game!","link":"https://example.com","version":1}]` + "\n"
	assetJSONQueryExpected      = `[{"id":1,"name":"cheeseboiger","tags":null},{"id":2,"name":"asset two","tags":["cheeseboiger"]}]` + "\n"
	assetJSONQueryExpectedLimit = `[{"id":1,"name":"cheeseboiger","tags":null}]` + "\n"
	assetJSONUpdate             = `{"name":"Updated cool asset","version":1}`
	assetJSONUpdateExpected     = `{"id":1,"name":"Updated cool asset","description":"A very nice asset to use in your game!","link":"https://example.com","version":2}` + "\n"
)

func (r *mockAssetRepo) CreateAsset(asset models.Asset) (*models.Asset, error) {
	id := 1
	asset.ID = &id
	asset.Version = &id
	r.data[id] = asset
	return &asset, nil
}
func (r *mockAssetRepo) FindAssetByID(id int) (*models.Asset, error) {
	a, ok := r.data[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &a, nil
}
func (r *mockAssetRepo) FindAssets(query []string, limit uint64) (*[]models.Asset, error) {
	var (
		assetIDs    []int          = []int{1, 2, 3}
		assetTitles []string       = []string{"cheeseboiger", "asset two", "asset three"}
		assetTags   [][]string     = [][]string{nil, []string{"cheeseboiger"}, nil}
		assets      []models.Asset = []models.Asset{
			{ID: &assetIDs[0], Name: &assetTitles[0], Tags: &assetTags[0]},
			{ID: &assetIDs[1], Name: &assetTitles[1], Tags: &assetTags[1]},
			{ID: &assetIDs[2], Name: &assetTitles[2], Tags: &assetTags[2]},
		}
		resultAssets []models.Asset
	)

	if len(query) != 0 {
		for _, a := range assets {
			if *a.Name == query[0] {
				resultAssets = append(resultAssets, a)
			}
			if slices.Contains(*a.Tags, query[0]) {
				resultAssets = append(resultAssets, a)
			}
		}
	} else {
		for _, v := range r.data {
			resultAssets = append(resultAssets, v)
		}
	}
	if limit != 0 {
		resultAssets = resultAssets[:limit]
	}
	return &resultAssets, nil
}
func (r *mockAssetRepo) UpdateAsset(id int, asset models.Asset) (*models.Asset, error) {
	var a models.Asset
	_, ok := r.data[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	a = r.data[id]
	if *a.Version != *asset.Version {
		return nil, r.ConflictErr()
	}
	if asset.Name != nil {
		a.Name = asset.Name
		n := *asset.Version + 1
		a.Version = &n
		r.data[id] = a
	}
	a = r.data[id]
	return &a, nil
}
func (r *mockAssetRepo) DeleteAsset(id int) error {
	_, ok := r.data[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.data, id)
	return nil
}
func (r *mockAssetRepo) NotFoundErr() error { return r.notFoundErr }
func (r *mockAssetRepo) ConflictErr() error { return r.conflictErr }

func TestPostAsset(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/assets", strings.NewReader(assetJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.PostAsset(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, assetJSONExpected, rec.Body.String())
	}
}

func TestGetAssetByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssetById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, assetJSONExpected, rec.Body.String())
	}
}

func TestGetAssetByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("4")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssetById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestGetAssets(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssets(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, assetJSONExpectedMany, rec.Body.String())
	}
}

func TestPatchAsset(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/assets", strings.NewReader(assetJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.PatchAsset(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, assetJSONUpdateExpected, rec.Body.String())
	}
}

func TestGetAssetsQuery(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/assets?q=%v", query), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssets(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, assetJSONQueryExpected, rec.Body.String())
	}
}

func TestGetAssetsQueryLimit(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/game-hangar/v1/assets?q=%v&l=%v", query, queryLimit), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssets(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, assetJSONQueryExpectedLimit, rec.Body.String())
	}
}

func TestPatchAssetNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/assets", strings.NewReader(assetJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("4")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.PatchAsset(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPatchAssetConflict(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/assets", strings.NewReader(assetJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.PatchAsset(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, conflictResponse, rec.Body.String())
	}
}

func TestDeleteAsset(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssetById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteAssetNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/assets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("5")
	h := &AssetHandler{logger: e.Logger, validator: v, repository: &ma}

	// Assertions
	if assert.NoError(t, h.GetAssetById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}
