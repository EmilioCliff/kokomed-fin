-- name: CreateUser :execresult
INSERT INTO users (full_name, phone_number, email, password, refresh_token, role, branch_id, updated_at, updated_by, created_by) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: UpdateUserPassword :execresult
UPDATE users SET password = ?, password_updated = password_updated + 1 WHERE email = ?;

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

-- name: HelperUser :many
SELECT id, full_name FROM users;

-- name: HelperUserById :one
SELECT full_name FROM users
WHERE id = ?;

-- name: ListUsersByCategory :many
SELECT 
    u.*, 
    b.name AS branch_name
FROM users u
JOIN branches b ON u.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(u.full_name) LIKE ?
        OR LOWER(u.email) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(u.role, ?) > 0
    )
 ORDER BY u.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsersByCategory :one
SELECT COUNT(*) AS total_loans
FROM users u
JOIN branches b ON u.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(u.full_name) LIKE ?
        OR LOWER(u.email) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(u.role, ?) > 0
    );

-- name: CheckUserExistance :one
SELECT COUNT(*) AS user_count FROM users WHERE email = ? LIMIT 1;