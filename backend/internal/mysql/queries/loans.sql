-- name: CreateLoan :execresult
INSERT INTO loans (product_id, client_id, loan_officer, loan_purpose, due_date, approved_by, disbursed_on, disbursed_by, total_installments, installments_period, status, processing_fee, fee_paid, created_by) 
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
    sqlc.arg("fee_paid"),
    sqlc.arg("created_by")
);

-- name: DisburseLoan :execresult
UPDATE loans 
    SET disbursed_on = ?,
    disbursed_by = ?,
    status = ?,
    due_date = ?
WHERE id = ?;

-- name: UpdateLoan :execresult
UPDATE loans 
    SET paid_amount = paid_amount + sqlc.arg("paid_amount"),
    updated_by = coalesce(sqlc.arg("updated_by"), updated_by)
WHERE id = sqlc.arg("id");

-- name: TransferLoan :execresult
UPDATE loans SET loan_officer = ?, updated_by = ? WHERE id = ?;

-- name: GetLoan :one
SELECT * FROM loans WHERE id = ? LIMIT 1;

-- name: GetClientActiveLoan :one
SELECT id FROM loans WHERE client_id = ? AND status = ? LIMIT 1;

-- name: ListLoans :many
SELECT 
    l.*, 
    p.branch_id,
    c.full_name AS client_name,
    u.full_name AS loan_officer_name
FROM loans l
JOIN products p ON l.product_id = p.id
JOIN clients c ON l.client_id = c.id
JOIN users u ON l.loan_officer = u.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(u.full_name) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(l.status, ?) > 0
    )
LIMIT ? OFFSET ?;


-- name: CountLoans :one
SELECT COUNT(*) AS total_loans 
FROM loans l
JOIN products p ON l.product_id = p.id
JOIN clients c ON l.client_id = c.id
JOIN users u ON l.loan_officer = u.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(u.full_name) LIKE ?
    )
    AND (
        COALESCE(?, '') = '' 
        OR FIND_IN_SET(l.status, ?) > 0
    );

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

-- name: UpdateLoanProcessingFeeStatus :execresult
UPDATE loans SET fee_paid = ? WHERE id = ?;

-- name: CheckActiveLoanForClient :one
SELECT EXISTS (
    SELECT 1
    FROM loans
    WHERE client_id = ? AND status = 'ACTIVE'
) AS has_active_loan LIMIT 1;

-- name: GetLoanPaymentData :one
SELECT 
    l.id AS loan_id,
    l.client_id,
    l.processing_fee,
    l.fee_paid,
    l.paid_amount,
    c.phone_number,
    p.repay_amount
FROM loans l
JOIN products p ON l.product_id = p.id
JOIN clients c ON l.client_id = c.id
WHERE l.id = ?
LIMIT 1;
