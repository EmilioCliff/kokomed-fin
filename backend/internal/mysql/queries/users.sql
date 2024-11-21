-- name: CreateUser :execresult
INSERT INTO users (full_name, phone_number, email, password, refresh_token, role, branch_id, updated_at, updated_by, created_by) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT password FROM users WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY full_name DESC;

-- name: UpdateUser :execresult
UPDATE users 
    SET full_name = ?, 
    phone_number = ?, 
    updated_at = ?, 
    updated_by = ? 
WHERE id = ?;

-- name: UpdateUserRole :execresult
UPDATE users
    SET role = ?,
    updated_at = ?, 
    updated_by = ? 
WHERE id = ?;

-- name: UpdateUserRefreshToken :execresult
UPDATE users
    SET refresh_token = ?
WHERE id = ?;

-- name: UpdateUserPassword :execresult
UPDATE users
    SET password = ?
WHERE id = ?;

-- name: UpdateUserBranch :execresult
UPDATE users
    SET branch_id = ?,
    updated_at = ?, 
    updated_by = ? 
WHERE id = ?;