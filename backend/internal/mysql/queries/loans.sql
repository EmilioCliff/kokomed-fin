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
    SET disbursed_on = coalesce(sqlc.narg("disbursed_on"), disbursed_on),
    disbursed_by = sqlc.arg("disbursed_by"),
    status = coalesce(sqlc.narg("status"), status),
    due_date = coalesce(sqlc.narg("due_date"), due_date),
    fee_paid =coalesce(sqlc.narg("fee_paid"), fee_paid)
WHERE id = sqlc.arg("id");

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

-- name: GetLoanData :many
SELECT id FROM loans;

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
 ORDER BY l.created_at DESC
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

-- name: ListExpectedPayments :many
SELECT 
	b.name AS branch_name,
	c.full_name AS client_name,
	u.full_name AS loan_officer_name,
	l.id AS loan_id, 
	p.loan_amount,
	p.repay_amount,
	COALESCE(p.repay_amount - l.paid_amount, 0) AS total_unpaid, 
	l.due_date
FROM clients c
JOIN loans l ON l.client_id = c.id AND l.status = 'ACTIVE'
JOIN products p ON l.product_id = p.id
JOIN users u ON l.loan_officer = u.id
JOIN branches b ON u.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(u.full_name) LIKE ?
    )
ORDER BY l.due_date DESC
LIMIT ? OFFSET ?;

-- name: CountExpectedPayments :one
SELECT COUNT(*) AS total_unexpected
FROM clients c
JOIN loans l ON l.client_id = c.id AND l.status = 'ACTIVE'
JOIN products p ON l.product_id = p.id
JOIN users u ON l.loan_officer = u.id
JOIN branches b ON u.branch_id = b.id
WHERE 
    (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(u.full_name) LIKE ?
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

-- name: GetLoanEvents :many
SELECT 
    l.id AS loan_id,
    l.disbursed_on AS disbursed_date,
    CASE 
        WHEN l.status = 'ACTIVE' THEN l.due_date
        ELSE NULL
    END AS due_date,

    c.full_name AS client_name,
    p.loan_amount AS loan_amount,
    CASE 
        WHEN l.status = 'ACTIVE' THEN (p.repay_amount - l.paid_amount)
        ELSE NULL
    END AS payment_due

FROM loans l
JOIN clients c ON l.client_id = c.id
JOIN products p ON l.product_id = p.id
WHERE l.disbursed_on IS NOT NULL 
ORDER BY l.disbursed_on DESC;

-- name: GetActiveLoanDetails :one
SELECT 
    l.id,
    p.loan_amount,
    p.repay_amount,
    l.disbursed_on,
    l.due_date,
    l.paid_amount
FROM loans l
JOIN products p ON l.product_id = p.id
WHERE 
    l.client_id = ? 
    AND l.status = 'ACTIVE';
