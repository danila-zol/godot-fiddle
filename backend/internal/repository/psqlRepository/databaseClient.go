package psqlRepository

import "github.com/jackc/pgx/v5/pgxpool"

type psqlDatabaseClient interface {
	AcquireConn() (*pgxpool.Conn, error)
}
