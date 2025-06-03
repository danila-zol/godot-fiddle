package psqlCasbinClient

import (
	"os"
	"slices"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	testCasbinClient *CasbinClient
	role             []string = []string{"Sharif", "Jet Fighters", "Call in"}

	object       string     = "Cadillac Deville"
	action       string     = "Fill with gas"
	roleExtented [][]string = [][]string{
		{"Sheikh", object, "Cruise"},
		{"Sheikh", object, action},
	}
)

func init() {
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../../.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	testCasbinClient, err = CasbinConfig{}.NewCasbinClient(
		os.Getenv("PSQL_CONNSTRING"),
		wd+"/rbac_model.conf",
	)
	if err != nil {
		panic("Error creating Casbin enforcer: " + err.Error())
	}
}

func TestAddPermissions(t *testing.T) {
	_, err := testCasbinClient.AddPermissions(role)
	assert.NoError(t, err)
	objs, err := testCasbinClient.enforcer.GetAllObjects()
	assert.NoError(t, err)
	acts, err := testCasbinClient.enforcer.GetAllActions()
	if assert.NoError(t, err) {
		assert.True(t, slices.Contains(objs, role[1]))
		assert.True(t, slices.Contains(acts, role[2]))
	}
}

func TestEnforcePass(t *testing.T) {
	pass, err := testCasbinClient.Enforce(role[0], role[1], role[2])
	assert.NoError(t, err)
	assert.True(t, pass)
}

func TestRemovePermissions(t *testing.T) {
	changed, err := testCasbinClient.RemovePermissions(role)
	assert.NoError(t, err)
	assert.True(t, changed)
	objs, err := testCasbinClient.enforcer.GetAllObjects()
	if assert.NoError(t, err) {
		assert.False(t, slices.Contains(objs, role[1]))
	}
}

func TestEnforceFail(t *testing.T) {
	pass, err := testCasbinClient.Enforce(role[0], role[1], role[2])
	assert.NoError(t, err)
	assert.False(t, pass)
}

func TestRemovePermissionsForObject(t *testing.T) {
	_, err := testCasbinClient.AddPermissions(roleExtented[0])
	assert.NoError(t, err)
	_, err = testCasbinClient.AddPermissions(roleExtented[1])
	assert.NoError(t, err)

	_, err = testCasbinClient.RemovePermissionsForObject(object, action)
	assert.NoError(t, err)

	actions, err := testCasbinClient.enforcer.GetAllActions()
	if assert.NoError(t, err) {
		assert.False(t, slices.Contains(actions, action))
		assert.True(t, slices.Contains(actions, "Cruise"))
	}
	_, err = testCasbinClient.RemovePermissionsForObject(object, roleExtented[0][2])
}
