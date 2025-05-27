package psqlRepository

import (
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
)

type psqlDatabaseClient interface {
	AcquireConn() (*pgxpool.Conn, error)
	ErrNoRows() error
}

type Enforcer interface {
	AddPermissions(...any) (bool, error)
	RemovePermissions(...any) (bool, error)
	RemovePermissionsForObject(obj, act string) (bool, error)
}

type ObjectUploader interface {
	PutObject(objectKey string, file io.Reader) error
	GetObjectLink(objectKey string) (*string, error)
	DeleteObject(objectKey string) error
}
