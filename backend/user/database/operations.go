package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TODO: MapNameToUser + assert one unique user per username!

func CreateUser(user User) (*User, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO "user".users
			(id, username, display_name, email, password, role_id, created_at, karma) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			(id, username, display_name, email, password, role_id, created_at, karma)`,
		user.ID, user.Username, user.Display_name, user.Email, user.Password,
		user.Role_id, user.Created_at, user.Karma,
	)

	err = row.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindFirstUser(id string) (*User, error) {
	var user User
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM "user".users WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Display_name, &user.Email, &user.Password,
		&user.Role_id, &user.Created_at, &user.Karma)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUsers() (*[]User, error) {
	var users []User

	conn, err := dbpool.Acquire(context.Background())
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
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Display_name, &user.Email,
			&user.Password, &user.Role_id, &user.Created_at, &user.Karma)
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

func UpdateUser(user User) (*User, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE "user".users SET 
		username=$1, display_name=$2, email=$3, password=$4, role_id=$5, created_at=$6, karma=$7
		WHERE id = $8`,
		user.Username, user.Display_name, user.Email, user.Password,
		user.Role_id, user.Created_at, user.Karma, user.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".users WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}

func CreateRole(role Role) (*Role, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO "user".roles
		(id, name) 
		VALUES
		($1, $2)
		RETURNING
		(id, name)`,
		role.ID, role.Name,
	)

	err = row.Scan(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FindFirstRole(id string) (*Role, error) {
	var role Role
	conn, err := dbpool.Acquire(context.Background())
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

func FindRoles() (*[]Role, error) {
	var roles []Role

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT * FROM "user".roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func UpdateRole(role Role) (*Role, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE "user".roles SET 
		name=$1
		WHERE id = $2`,
		role.Name, role.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func DeleteRole(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM "user".roles WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}
