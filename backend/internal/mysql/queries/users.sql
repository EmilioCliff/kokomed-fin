-- name: CreateUser :execresult
INSERT INTO users (full_name, phone_number, email, password, refresh_token, role, branch_id, updated_at, updated_by, created_by) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT password FROM users WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY full_name DESC LIMIT ? OFFSET ?;

-- name: UpdateUser :execresult
UPDATE users 
    SET role = coalesce(sqlc.narg("role"), role), 
    branch_id = coalesce(sqlc.narg("branch_id"), branch_id),
    password = coalesce(sqlc.narg("password"), password),
    refresh_token = coalesce(sqlc.narg("refresh_token"), refresh_token),
    updated_at = coalesce(sqlc.narg("updated_at"), updated_at), 
    updated_by = coalesce(sqlc.narg("updated_by"), updated_by) 
WHERE id = sqlc.arg("id");