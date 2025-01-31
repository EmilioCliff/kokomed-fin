-- name: ListBranches :many
SELECT * FROM branches;

-- name: GetBranch :one
SELECT * FROM branches WHERE id = ? LIMIT 1;

-- name: CreateBranch :execresult
INSERT INTO branches (name) VALUES (?);

-- name: UpdateBranch :execresult
UPDATE branches SET name = ? WHERE id = ?;

-- name: DeleteBranch :exec
DELETE FROM branches WHERE id = ?;

-- name: ListBrachesByCategory :many
SELECT * FROM branches b
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(b.name) LIKE ?
    )
LIMIT ? OFFSET ?;

-- name: CountBranchesByCategory :one
SELECT COUNT(*) AS total_branches
FROM branches b
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(b.name) LIKE ?
    )