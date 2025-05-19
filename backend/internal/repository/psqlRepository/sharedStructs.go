package psqlRepository

import "github.com/jackc/pgx/v5/pgxpool"

type psqlDatabaseClient interface {
	AcquireConn() (*pgxpool.Conn, error)
	ErrNoRows() error
}

type Enforcer interface {
	AddPermissions(...any) (bool, error)
	RemovePermissions(...any) (bool, error)
	RemovePermissionsForObject(obj, act string) (bool, error)
}
