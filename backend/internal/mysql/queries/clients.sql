-- name: CreateClient :execresult
INSERT INTO clients (full_name, phone_number, id_number, dob, gender, active, branch_id, assigned_staff, updated_by, updated_at, created_by) 
VALUES (
    sqlc.arg("full_name"),
    sqlc.arg("phone_number"),
    sqlc.narg("id_number"),
    sqlc.narg("dob"),
    sqlc.arg("gender"),
    sqlc.arg("active"),
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
    SET full_name = coalesce(sqlc.narg("full_name"), full_name),
    phone_number = coalesce(sqlc.narg("phone_number"), phone_number),
    id_number = coalesce(sqlc.narg("id_number"), id_number),
    dob = coalesce(sqlc.narg("dob"), dob),
    gender = coalesce(sqlc.narg("gender"), gender),
    active = coalesce(sqlc.narg("active"), active),
    branch_id = coalesce(sqlc.narg("branch_id"), branch_id),
    assigned_staff = coalesce(sqlc.narg("assigned_staff"), assigned_staff),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = sqlc.arg("updated_by")
WHERE id = sqlc.arg("id");

-- name: DeleteClient :execresult
DELETE FROM clients WHERE id = ?;

-- name: GetClientByPhoneNumber :one
SELECT * FROM clients WHERE phone_number = ? LIMIT 1;

-- name: ListClientsByBranch :many
SELECT * FROM clients WHERE branch_id = ? LIMIT ? OFFSET ?;

-- name: ListClientsByActiveStatus :many
SELECT * FROM clients WHERE active = ? LIMIT ? OFFSET ?;

-- name: UpdateClientOverpayment :execresult
UPDATE clients SET overpayment = ? WHERE id = ?;
