// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (full_name, phone_number, email, password, refresh_token, role, branch_id, updated_at, updated_by, created_by) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	Role         UsersRole `json:"role"`
	BranchID     uint32    `json:"branch_id"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    uint32    `json:"updated_by"`
	CreatedBy    uint32    `json:"created_by"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.FullName,
		arg.PhoneNumber,
		arg.Email,
		arg.Password,
		arg.RefreshToken,
		arg.Role,
		arg.BranchID,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.CreatedBy,
	)
}

const getUser = `-- name: GetUser :one
SELECT id, full_name, phone_number, email, password, refresh_token, role, branch_id, updated_by, updated_at, created_by, created_at FROM users WHERE id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uint32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.PhoneNumber,
		&i.Email,
		&i.Password,
		&i.RefreshToken,
		&i.Role,
		&i.BranchID,
		&i.UpdatedBy,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT password FROM users WHERE email = ? LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var password string
	err := row.Scan(&password)
	return password, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, full_name, phone_number, email, password, refresh_token, role, branch_id, updated_by, updated_at, created_by, created_at FROM users ORDER BY full_name DESC LIMIT ? OFFSET ?
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.PhoneNumber,
			&i.Email,
			&i.Password,
			&i.RefreshToken,
			&i.Role,
			&i.BranchID,
			&i.UpdatedBy,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :execresult
UPDATE users 
    SET role = coalesce(?, role), 
    branch_id = coalesce(?, branch_id),
    password = coalesce(?, password),
    refresh_token = coalesce(?, refresh_token),
    updated_at = coalesce(?, updated_at), 
    updated_by = coalesce(?, updated_by) 
WHERE id = ?
`

type UpdateUserParams struct {
	Role         NullUsersRole  `json:"role"`
	BranchID     sql.NullInt32  `json:"branch_id"`
	Password     sql.NullString `json:"password"`
	RefreshToken sql.NullString `json:"refresh_token"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	UpdatedBy    sql.NullInt32  `json:"updated_by"`
	ID           uint32         `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Role,
		arg.BranchID,
		arg.Password,
		arg.RefreshToken,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.ID,
	)
}