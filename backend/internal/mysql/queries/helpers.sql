-- name: DashBoardRecentsPayments :many
SELECT id, paying_name, amount, paid_date 
FROM non_posted 
WHERE paid_date >= NOW() - INTERVAL 15 DAY;

-- name: DashBoardInactiveLoans :many
SELECT l.id, p.loan_amount, p.repay_amount, l.client_id, l.loan_officer, l.approved_by, l.created_at FROM loans l 
JOIN products p ON l.product_id = p.id WHERE l.status = 'INACTIVE';

-- name: DashBoardDataHelper :one
SELECT
    -- Clients
    (SELECT COUNT(*) FROM clients) AS total_clients,
    (SELECT COUNT(*) FROM clients WHERE active = TRUE) AS active_clients,

    -- Loans
    (SELECT COUNT(*) FROM loans) AS total_loans,
    (SELECT COUNT(*) FROM loans WHERE status = 'ACTIVE') AS active_loans,
    (SELECT COUNT(*) FROM loans WHERE status = 'INACTIVE') AS inactive_loans,

     -- Financials
    (SELECT COALESCE(SUM(p.loan_amount), 0) 
     FROM loans l 
     JOIN products p ON l.product_id = p.id) AS total_loan_amount,

    (SELECT COALESCE(SUM(p.loan_amount), 0) 
     FROM loans l 
     JOIN products p ON l.product_id = p.id 
     WHERE l.status != 'INACTIVE') AS total_loan_disbursed,

    (SELECT COALESCE(SUM(p.loan_amount), 0) 
     FROM loans l 
     JOIN products p ON l.product_id = p.id 
     WHERE l.status = 'COMPLETED') AS total_loan_paid,

    -- Non-posted Payments
    (SELECT COALESCE(SUM(amount), 0) FROM non_posted) AS total_payments_received,
    (SELECT COALESCE(SUM(amount), 0) FROM non_posted WHERE assign_to IS NULL) AS total_non_posted;