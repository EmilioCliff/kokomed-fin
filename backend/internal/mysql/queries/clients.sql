-- name: CreateClient :execresult
INSERT INTO clients (full_name, phone_number, id_number, dob, gender, active, branch_id, assigned_staff, updated_by, updated_at, created_by) 
VALUES (
    sqlc.arg("full_name"),
    sqlc.arg("phone_number"),
    sqlc.narg("id_number"),
    sqlc.narg("dob"),
    sqlc.arg("gender"),
    true,
    sqlc.arg("branch_id"),
    sqlc.arg("assigned_staff"),
    sqlc.arg("updated_by"),
    CURRENT_TIMESTAMP,
    sqlc.arg("created_by")
);

-- name: ListClients :many
SELECT * FROM clients LIMIT ? OFFSET ?;

-- name: GetClient :one
SELECT * FROM clients WHERE id = ? LIMIT 1;

-- name: UpdateClient :execresult
UPDATE clients 
    SET id_number = coalesce(sqlc.narg("id_number"), id_number),
    dob = coalesce(sqlc.narg("dob"), dob),
    active = coalesce(sqlc.narg("active"), active),
    branch_id = coalesce(sqlc.narg("branch_id"), branch_id),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = sqlc.arg("updated_by")
WHERE id = sqlc.arg("id");

-- -- name: UpdateClientOverpayment :execresult
-- UPDATE clients 
--     SET overpayment = overpayment + sqlc.arg("overpayment")
-- WHERE phone_number = sqlc.arg("phone_number");

-- name: UpdateClientOverpayment :execresult
UPDATE clients
SET overpayment = overpayment + sqlc.arg("overpayment")
WHERE 
    (phone_number = sqlc.arg("phone_number") AND sqlc.arg("phone_number") IS NOT NULL)
    OR 
    (id = sqlc.arg("client_id") AND sqlc.arg("client_id") IS NOT NULL);

-- name: NullifyClientOverpayment :execresult
UPDATE clients
SET overpayment = 0
WHERE id = ?;


-- name: DeleteClient :execresult
DELETE FROM clients WHERE id = ?;

-- name: GetClientIDByPhoneNumber :one
SELECT id FROM clients WHERE phone_number = ? LIMIT 1;

-- name: GetClientOverpayment :one
SELECT overpayment FROM clients WHERE id = ? LIMIT 1;

-- name: ListClientsByBranch :many
SELECT * FROM clients WHERE branch_id = ? LIMIT ? OFFSET ?;

-- name: ListClientsByActiveStatus :many
SELECT * FROM clients WHERE active = ? LIMIT ? OFFSET ?;

-- name: HelperClient :many
SELECT id, full_name, phone_number FROM clients;

-- name: ListClientsByCategory :many
SELECT 
    c.id, c.full_name, c.phone_number, c.id_number, c.dob, c.gender, c.active, 
    c.branch_id, c.assigned_staff, c.overpayment, c.updated_by, 
    c.updated_at, c.created_at, c.created_by, 
    b.name AS branch_name,
    COALESCE(SUM(DISTINCT COALESCE(p.loan_amount, 0) - COALESCE(l.paid_amount, 0)), 0) AS dueAmount
FROM clients c
JOIN branches b ON c.branch_id = b.id
LEFT JOIN loans l ON c.id = l.client_id AND l.status = 'ACTIVE'
LEFT JOIN products p ON l.product_id = p.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(c.phone_number) LIKE ?
    )
    AND (
        sqlc.narg('active') IS NULL OR c.active = sqlc.narg('active')
    )
GROUP BY 
    c.id, c.full_name, c.phone_number, c.id_number, c.dob, c.gender, c.active, 
    c.branch_id, c.assigned_staff, c.overpayment, c.updated_by, 
    c.updated_at, c.created_at, c.created_by, b.name
ORDER BY c.created_at DESC
LIMIT ? OFFSET ?;


-- name: CountClientsByCategory :one
SELECT COUNT(*) AS total_clients
FROM clients c
JOIN branches b ON c.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(c.phone_number) LIKE ?
    )
    AND (
        sqlc.narg('active') IS NULL OR c.active = sqlc.narg('active')
    );
