-- name: CreateLoan :execresult
INSERT INTO loans (product_id, client_id, loan_officer, loan_purpose, due_date, approved_by, disbursed_on, disbursed_by, total_installments, installments_period, status, processing_fee, created_by) 
VALUES (
    sqlc.arg("product_id"),
    sqlc.arg("client_id"),
    sqlc.arg("loan_officer"),
    sqlc.narg("loan_purpose"),
    sqlc.narg("due_date"),
    sqlc.arg("approved_by"),
    sqlc.narg("disbursed_on"),
    sqlc.narg("disbursed_by"),
    sqlc.arg("total_installments"),
    sqlc.arg("installments_period"),
    sqlc.arg("status"),
    sqlc.arg("processing_fee"),
    sqlc.arg("created_by")
);

-- name: DisburseLoan :execresult
UPDATE loans 
    SET disbursed_on = ?,
    disbursed_by = ?,
    due_date = ?
WHERE id = ?;

-- name: UpdateLoan :execresult
UPDATE loans 
    SET paid_amount = sqlc.arg("paid_amount"),
    updated_by = coalesce(sqlc.arg("updated_by"), updated_by)
WHERE id = sqlc.arg("id");

-- name: TransferLoan :execresult
UPDATE loans SET loan_officer = ? WHERE id = ?;

-- name: GetLoan :one
SELECT * FROM loans WHERE id = ? LIMIT 1;

-- name: ListLoans :many
SELECT * FROM loans LIMIT ? OFFSET ?;

-- name: DeleteLoan :exec
DELETE FROM loans WHERE id = ?;

-- name: ListLoansByClient :many
SELECT * FROM loans WHERE client_id = ? LIMIT ? OFFSET ?;

-- name: ListLoansByLoanOfficer :many
SELECT * FROM loans WHERE loan_officer = ? LIMIT ? OFFSET ?;

-- name: ListLoansByStatus :many
SELECT * FROM loans WHERE status = ? LIMIT ? OFFSET ?;

-- name: ListNonDisbursedLoans :many
SELECT * FROM loans WHERE disbursed_on IS NULL LIMIT ? OFFSET ?;

-- name: UpdateLoanStatus :execresult
UPDATE loans SET status = ? WHERE id = ?;
