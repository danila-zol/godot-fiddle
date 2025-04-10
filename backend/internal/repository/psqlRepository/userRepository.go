package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
)

type PsqlUserRepository struct {
	databaseClient *psqlDatabase.PsqlDatabaseClient
	notFoundErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func (pur *PsqlUserRepository) NewPsqlUserRepository(dbClient *psqlDatabase.PsqlDatabaseClient) (*PsqlUserRepository, error) {
	return &PsqlUserRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
	}, nil
}

// TODO: MapNameToUser + assert one unique user per username!

func (pur *PsqlUserRepository) CreateUser(user models.User) (*models.User, error) {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".users
			(id, username, displayName, email, password, roleID, createdAt, karma) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			(id, username, displayName, email, password, roleID, createdAt, karma)`,
		user.ID, user.Username, user.DisplayName, user.Email, user.Password,
		user.RoleID, user.CreatedAt, user.Karma,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pur *PsqlUserRepository) FindFirstUser(id string) (*models.User, error) {
	var user models.User
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM "user".users WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.Password,
		&user.RoleID, &user.CreatedAt, &user.Karma)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pur *PsqlUserRepository) FindUsers() (*[]models.User, error) {
	var users []models.User

	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT * FROM "user".users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email,
			&user.Password, &user.RoleID, &user.CreatedAt, &user.Karma)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (pur *PsqlUserRepository) UpdateUser(id string, user models.User) (*models.User, error) {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE "user".users SET 
		username=$1, display_name=$2, email=$3, password=$4, role_id=$5, created_at=$6, karma=$7
		WHERE id = $8
		RETURNING
		(id, username, displayName, email, password, roleID, createdAt, karma)`,
		user.Username, user.DisplayName, user.Email, user.Password,
		user.RoleID, user.CreatedAt, user.Karma, id,
	).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pur *PsqlUserRepository) DeleteUser(id string) error {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".users WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pur.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pur *PsqlUserRepository) CreateRole(role models.Role) (*models.Role, error) {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".roles
		(id, name) 
		VALUES
		($1, $2)
		RETURNING
		(id, name)`,
		role.ID, role.Name,
	).Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (pur *PsqlUserRepository) FindRoleByID(id string) (*models.Role, error) {
	var role models.Role
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM "user".roles WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (pur *PsqlUserRepository) UpdateRole(role models.Role) (*models.Role, error) {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE "user".roles SET 
		name=$1
		WHERE id = $2
		RETURNING
		(id, name)`,
		role.Name, role.ID,
	).Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (pur *PsqlUserRepository) DeleteRole(id string) error {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".roles WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pur.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pur *PsqlUserRepository) CreateSession(session models.Session) (*models.Session, error) {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO "user".sessions
		(id, userID) 
		VALUES
		($1, $2)
		RETURNING
		(id, userID)`,
		session.ID, session.UserID,
	).Scan(&session.ID, &session.UserID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (pur *PsqlUserRepository) FindSessionByID(id string) (*models.Session, error) {
	var session models.Session
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM "user".sessions WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&session.ID, &session.UserID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (pur *PsqlUserRepository) DeleteSession(id string) error {
	conn, err := pur.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".sessions WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pur.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}
