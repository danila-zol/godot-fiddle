package psqlCasbinClient

import (
	"github.com/casbin/casbin/v2"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v3"
)

type CasbinConfig struct{}

type CasbinClient struct {
	enforcer *casbin.Enforcer
}

func (c CasbinConfig) NewCasbinClient(connstring, modelPath string) (*CasbinClient, error) {
	a, _ := pgxadapter.NewAdapter(connstring, pgxadapter.WithDatabase("gamehangar"))

	ce, err := casbin.NewEnforcer(modelPath, a)
	if err != nil {
		return nil, err
	}

	_, err = ce.AddPolicies([][]string{
		{"admin"},
		{"freetier", "assets", "POST"},
		{"freetier", "demos", "POST"},
		{"freetier", "threads", "POST"},
		{"freetier", "messages", "POST"},
		{"paidtier", "demos", "POSTExtended"},
		{"paidtier", "freetier"},
	})
	if err != nil {
		return nil, err
	}
	ce.LoadPolicy()

	return &CasbinClient{enforcer: ce}, nil
}

func (c *CasbinClient) AddPermissions(params ...any) (bool, error) {
	return c.enforcer.AddPolicy(params...)
}

func (c *CasbinClient) RemovePermissions(params ...any) (bool, error) {
	return c.enforcer.RemovePolicy(params...)
}

func (c *CasbinClient) RemovePermissionsForObject(obj, act string) (bool, error) {
	return c.enforcer.RemoveFilteredPolicy(0, "", obj, act)
}

func (c *CasbinClient) Enforce(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}
