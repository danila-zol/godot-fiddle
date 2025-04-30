package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/domain/models"
)

type PsqlUserRepository struct {
	databaseClient psqlDatabaseClient
	conflictErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlUserRepository(dbClient psqlDatabaseClient) *PsqlUserRepository {
	return &PsqlUserRepository{
		databaseClient: dbClient,
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
			(username, display_name, email, password, role_id) 
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			(id, username, display_name, email, password, verified, role_id, created_at, karma)`,
		user.Username, user.DisplayName, user.Email, user.Password, user.RoleID,
	).Scan(&user) // Mind the field order!
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUserByID(id string) (*models.User, error) {
	var user models.User
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, username, display_name, email, password, verified, role_id, created_at, karma)
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
		`SELECT (id, username, display_name, email, password, verified, role_id, created_at, karma) 
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
		`SELECT (id, username, display_name, email, password, verified, role_id, created_at, karma)
		FROM "user".users WHERE username = $1 LIMIT 1`,
		username,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) FindUsers() (*[]models.User, error) {
	var users []models.User

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, username, display_name, email, password, verified, role_id, created_at, karma)
		FROM "user".users`,
	)
	if err != nil {
		return nil, err
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

func (r *PsqlUserRepository) UpdateUser(id string, user models.User) (*models.User, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE "user".users SET 
			username=COALESCE($1, username), display_name=COALESCE($2, display_name), email=COALESCE($3, email), 
			password=COALESCE($4, password), verified=COALESCE($5, verified), role_id=COALESCE($6, role_id), 
			karma=COALESCE($7, karma)
		WHERE id = $8
		RETURNING
			(id, username, display_name, email, password, verified, role_id, created_at, karma)`,
		user.Username, user.DisplayName, user.Email, user.Password,
		user.Verified, user.RoleID, user.Karma, id,
	).Scan(&user) // Mind the field order!
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PsqlUserRepository) DeleteUser(id string) error {
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

func (r *PsqlUserRepository) CreateRole(role models.Role) (*models.Role, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".roles
		(name) 
		VALUES
		($1)
		RETURNING
		(id, name, version)`,
		role.Name,
	).Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PsqlUserRepository) FindRoleByID(id string) (*models.Role, error) {
	var role models.Role
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, name, version) FROM "user".roles WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PsqlUserRepository) UpdateRole(id string, role models.Role) (*models.Role, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(
		context.Background(),
		`SELECT (id) FROM "user".roles WHERE id=$1 AND version=$2`,
		id, *role.Version,
	).Scan(&id)
	if err != nil {
		return nil, r.ConflictErr()
	}

	err = conn.QueryRow(context.Background(),
		`UPDATE "user".roles SET 
		name=COALESCE($1, name)
		WHERE id = $2
		RETURNING
		(id, name, version)`,
		role.Name, id,
	).Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PsqlUserRepository) DeleteRole(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".roles WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
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

func (r *PsqlUserRepository) FindSessionByID(id string) (*models.Session, error) {
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

func (r *PsqlUserRepository) DeleteAllUserSessions(userID string) error {
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

func (r *PsqlUserRepository) DeleteSession(id string) error {
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
