-- name: CreateInstallment :execresult
INSERT INTO installments (loan_id, installment_number, amount_due, remaining_amount, due_date) 
VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetInstallment :one
SELECT * FROM installments WHERE id = ? LIMIT 1;

-- name: ListInstallmentsByLoan :many
SELECT * FROM installments WHERE loan_id = ? ORDER BY due_date ASC;

-- name: ListUnpaidInstallmentsByLoan :many
SELECT * FROM installments WHERE loan_id = ? AND remaining_amount > 0 ORDER BY due_date ASC;

-- name: UpdateInstallment :execresult
UPDATE installments 
    SET remaining_amount =  sqlc.arg("remaining_amount"),
    paid =  coalesce(sqlc.narg("paid"), paid),
    paid_at =  coalesce(sqlc.narg("paid_at"), paid_at)
WHERE id = sqlc.arg("id");

-- name: PayInstallment :execresult
UPDATE installments 
    SET remaining_amount = sqlc.arg("remaining_amount"),
    paid =  coalesce(sqlc.narg("paid"), paid),
    paid_at =  coalesce(sqlc.narg("paid_at"), paid_at)
WHERE id = sqlc.arg("id");

-- name: GetUnpaidInstallmentsData :many
SELECT 
    i.installment_number,
    i.amount_due,
    i.remaining_amount,
    i.due_date,

    l.id AS loan_id,
    p.loan_amount,
    p.repay_amount,
    l.paid_amount AS total_paid_amount,
    
    c.id AS client_id,
    c.full_name AS client_name,
    c.phone_number AS client_phone,
    b.name AS branch_name,

    (
        SELECT SUM(i2.remaining_amount)
        FROM installments i2
        WHERE i2.loan_id = l.id
          AND (i2.paid = FALSE OR i2.remaining_amount > 0)
          AND i2.due_date <= CURDATE()
    ) AS total_due_amount

FROM installments i
JOIN loans l ON i.loan_id = l.id
JOIN clients c ON l.client_id = c.id
JOIN branches b ON c.branch_id = b.id
JOIN products p ON l.product_id = p.id

WHERE 
    (i.paid = FALSE OR i.remaining_amount > 0) 
    AND i.due_date <= CURDATE()
    AND (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(c.phone_number) LIKE ?
    )

ORDER BY i.due_date DESC
LIMIT ? OFFSET ?;

-- name: CountUnpaidInstallmentsData :one
SELECT COUNT(*) AS total_unpaid_installments
FROM installments i
JOIN loans l ON i.loan_id = l.id
JOIN clients c ON l.client_id = c.id

WHERE 
    (i.paid = FALSE OR i.remaining_amount > 0) 
    AND i.due_date <= CURDATE()
    AND (
        COALESCE(?, '') = '' 
        OR LOWER(c.full_name) LIKE ?
        OR LOWER(c.phone_number) LIKE ?
    );
