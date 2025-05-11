package psqlRepository

import (
	"context"
	"errors"
	"fmt"
	"gamehangar/internal/domain/models"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PsqlUserRepository struct {
	databaseClient psqlDatabaseClient
	enforcer       Enforcer
	conflictErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlUserRepository(dbClient psqlDatabaseClient, e Enforcer) *PsqlUserRepository {
	return &PsqlUserRepository{
		databaseClient: dbClient,
		enforcer:       e,
		conflictErr:    errors.New("Record conflict!"),
	}
}
func (r *PsqlUserRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

// Returns "Record conflict!" to specify conflicting record versions on update
func (r *PsqlUserRepository) ConflictErr() error { return r.conflictErr }

func (r *PsqlUserRepository) CreateUser(user models.User) (*models.User, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".users
			(username, display_name, email, password, role) 
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			(id, username, display_name, email, password, verified, role, created_at, karma)`,
		user.Username, user.DisplayName, user.Email, user.Password, user.Role,
	).Scan(&user) // Mind the field order!
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, username, display_name, email, password, verified, role, created_at, karma)
		FROM "user".users WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, username, display_name, email, password, verified, role, created_at, karma) 
		FROM "user".users WHERE email = $1 LIMIT 1`,
		email,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, username, display_name, email, password, verified, role, created_at, karma)
		FROM "user".users WHERE username = $1 LIMIT 1`,
		username,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUsers(keywords []string, limit uint64) (*[]models.User, error) {
	var users []models.User

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query :=
			`SELECT (id, username, display_name, email, password, verified, role, created_at, karma) 
				FROM
					((SELECT id, username, display_name, email, password, verified, role, created_at, karma
					FROM "user".users
					WHERE asset_ts @@ to_tsquery_multilang($1))
				UNION
					(SELECT id, username, display_name, email, password, verified, role, created_at, karma 
					FROM "user".users
					WHERE tags && ($2) COLLATE case_insensitive))
			ORDER BY karma DESC`
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(),
			query, strings.Join(keywords, " | "), keywords,
		)
		if err != nil {
			return nil, err
		}
	} else {
		query := `SELECT 
			(id, username, display_name, email, password, verified, role, created_at, karma) 
			FROM "user".users
			ORDER BY karma DESC`
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, r.NotFoundErr()
	}
	return &users, nil
}

func (r *PsqlUserRepository) UpdateUser(id uuid.UUID, user models.User) (*models.User, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE "user".users SET 
			username=COALESCE($1, username), display_name=COALESCE($2, display_name), email=COALESCE($3, email), 
			password=COALESCE($4, password), verified=COALESCE($5, verified), role=COALESCE($6, role), 
			karma=COALESCE($7, karma)
		WHERE id = $8
		RETURNING
			(id, username, display_name, email, password, verified, role, created_at, karma)`,
		user.Username, user.DisplayName, user.Email, user.Password,
		user.Verified, user.Role, user.Karma, id,
	).Scan(&user) // Mind the field order!
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) DeleteUser(id uuid.UUID) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".users WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}

func (r *PsqlUserRepository) CreateRole(role string) error {
	_, err := r.enforcer.AddPermissions(role)
	return err
}

func (r *PsqlUserRepository) DeleteRole(role string) error {
	_, err := r.enforcer.RemovePermissions(role)
	return err
}

func (r *PsqlUserRepository) CreateSession(session models.Session) (*models.Session, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".sessions
		(user_id) 
		VALUES
		($1)
		RETURNING
		(id, user_id)`,
		session.UserID,
	).Scan(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PsqlUserRepository) FindSessionByID(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, user_id) FROM "user".sessions WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *PsqlUserRepository) DeleteAllUserSessions(userID uuid.UUID) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".sessions WHERE user_id=$1`, userID)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}

func (r *PsqlUserRepository) DeleteSession(id uuid.UUID) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".sessions WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}
