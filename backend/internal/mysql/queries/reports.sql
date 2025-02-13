-- name: GetPaymentReportData :many
SELECT 
    np.transaction_number,
    np.paying_name,
    np.amount,
    np.account_number,
    np.transaction_source,
    np.paid_date,
    COALESCE(c.full_name, 'Unassigned') AS assigned_name, 
    COALESCE(np.assigned_by, 'System') AS assigned_by 
FROM 
    non_posted np
LEFT JOIN 
    clients c ON np.assign_to = c.id
WHERE 
    np.paid_date BETWEEN ? AND ?;

-- name: GetBranchReportData :many
WITH branch_metrics AS (
    SELECT 
        b.id,
        b.name AS branch_name,
        COUNT(DISTINCT c.id) AS total_clients,
        COUNT(DISTINCT u.id) AS total_users
    FROM branches b
    LEFT JOIN clients c ON c.branch_id = b.id
    LEFT JOIN users u ON u.branch_id = b.id
    GROUP BY b.id, b.name
),
loan_metrics AS (
    SELECT 
        p.branch_id,
        SUM(CASE 
            WHEN l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") 
            THEN p.loan_amount 
            ELSE 0 
        END) AS total_disbursed_amount,
        SUM(CASE 
            WHEN l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") 
            THEN l.paid_amount 
            ELSE 0 
        END) AS total_collected_amount,
        SUM(CASE 
            WHEN l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") 
            THEN (p.repay_amount - l.paid_amount)
            ELSE 0 
        END) AS total_outstanding_amount,
        COUNT(DISTINCT CASE 
            WHEN l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") 
            THEN l.id 
        END) AS total_loans_issued,
        COUNT(DISTINCT CASE 
            WHEN l.status = 'DEFAULTED' 
            AND l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") 
            THEN l.id 
        END) AS defaulted_loans
    FROM products p
    JOIN loans l ON l.product_id = p.id 
    WHERE l.status != 'INACTIVE'
    GROUP BY p.branch_id
)
SELECT 
    bm.branch_name,
    bm.total_clients,
    bm.total_users,
    COALESCE(lm.total_loans_issued, 0) AS total_loans_issued,
    COALESCE(lm.total_disbursed_amount, 0) AS total_disbursed_amount,
    COALESCE(lm.total_collected_amount, 0) AS total_collected_amount,
    COALESCE(lm.total_outstanding_amount, 0) AS total_outstanding_amount,
    COALESCE(
        CASE 
            WHEN lm.total_loans_issued > 0 
            THEN (lm.defaulted_loans * 100.0) / lm.total_loans_issued 
            ELSE 0 
        END,
        0
    ) AS default_rate
FROM branch_metrics bm
LEFT JOIN loan_metrics lm ON lm.branch_id = bm.id
ORDER BY bm.branch_name;

-- name: GetUserAdminsReportData :many
SELECT 
    u.full_name AS name,
    u.role,
    b.name AS branch_name,
    COUNT(DISTINCT CASE 
        WHEN l.approved_by = u.id AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
    END) AS approved_loans,
    COUNT(DISTINCT CASE 
        WHEN l.loan_officer = u.id AND l.status = 'ACTIVE' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
    END) AS active_loans,

    COUNT(DISTINCT CASE 
        WHEN l.loan_officer = u.id AND l.status = 'COMPLETED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
    END) AS completed_loans,
    COALESCE(
        (COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id AND l.status = 'DEFAULTED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
        END) * 100.0) / NULLIF(COUNT(DISTINCT CASE 
            WHEN (l.loan_officer = u.id OR l.created_by = u.id OR l.approved_by = u.id OR l.disbursed_by = u.id OR l.updated_by = u.id) 
            AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
        END), 0), 
        0
    ) AS default_rate,
    COUNT(DISTINCT CASE 
        WHEN c.created_by = u.id AND c.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN c.id 
    END) AS clients_registered,
    COUNT(DISTINCT CASE 
        WHEN LOWER(np.assigned_by) = LOWER(u.full_name) AND np.paid_date BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN np.id 
    END) AS payments_assigned

FROM users u
LEFT JOIN branches b ON u.branch_id = b.id
LEFT JOIN loans l ON 
    l.loan_officer = u.id OR 
    l.created_by = u.id OR 
    l.approved_by = u.id OR 
    l.disbursed_by = u.id OR 
    l.updated_by = u.id
LEFT JOIN clients c ON c.created_by = u.id
LEFT JOIN non_posted np ON LOWER(np.assigned_by) = LOWER(u.full_name)
GROUP BY u.id, u.full_name, u.role, b.name
ORDER BY u.role DESC, u.full_name;

-- name: GetUserUsersReportData :one
SELECT 
    u.full_name AS name,
    u.role,
    b.name AS branch,

    COUNT(DISTINCT CASE 
        WHEN c.assigned_staff = u.id 
        AND c.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN c.id 
    END) AS total_clients_handled,

    COUNT(DISTINCT CASE 
        WHEN l.approved_by = u.id 
        AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
    END) AS loans_approved,

    COALESCE((
        SELECT SUM(p.repay_amount) 
        FROM loans l
        JOIN products p ON l.product_id = p.id
        WHERE l.loan_officer = u.id AND l.status != 'INACTIVE'
        AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")
    ), 0) AS total_loan_amount_managed,

    COALESCE((
        SELECT SUM(l.paid_amount) 
        FROM loans l
        WHERE l.loan_officer = u.id AND l.paid_amount > 0 
        AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")
    ), 0) AS total_collected_amount,

    COALESCE(
        (COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id 
            AND l.status = 'DEFAULTED' 
            AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
        END) * 100.0) 
        / NULLIF(
            COUNT(DISTINCT CASE 
                WHEN l.loan_officer = u.id AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id 
            END), 0
        ), 
        0
    ) AS default_rate,

    COUNT(DISTINCT CASE 
        WHEN LOWER(np.assigned_by) = LOWER(u.full_name) 
        AND np.paid_date BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN np.id 
    END) AS assigned_payments,

    (
        SELECT COALESCE(
            JSON_ARRAYAGG(
                JSON_OBJECT(
                    'loan_id', l.id,
                    'status', l.status,
                    'client_name', c.full_name,
                    'loan_amount', p.loan_amount,
                    'repay_amount', p.repay_amount,
                    'paid_amount', l.paid_amount,
                    'disbursed_on', l.disbursed_on
                )
            ), '[]'
        )
        FROM loans l
        JOIN clients c ON l.client_id = c.id
        JOIN products p ON l.product_id = p.id
        WHERE l.loan_officer = u.id 
    ) AS assigned_loans,

    (
        SELECT COALESCE(
            JSON_ARRAYAGG(
                JSON_OBJECT(
                    'transaction_number', np.transaction_number,
                    'client_name', assigned_client.full_name,
                    'amount_paid', np.amount,
                    'paid_date', np.paid_date
                )
            ), '[]'
        )
        FROM non_posted np
        LEFT JOIN clients assigned_client ON np.assign_to = assigned_client.id
        WHERE LOWER(np.assigned_by) = LOWER(u.full_name)
    ) AS assigned_payments_list

FROM users u
LEFT JOIN branches b ON u.branch_id = b.id
LEFT JOIN clients c ON c.assigned_staff = u.id
LEFT JOIN loans l ON 
    l.loan_officer = u.id 
    OR l.approved_by = u.id 
    OR l.created_by = u.id 
    OR l.disbursed_by = u.id 
    OR l.updated_by = u.id

LEFT JOIN products p ON l.product_id = p.id
LEFT JOIN non_posted np ON LOWER(np.assigned_by) = LOWER(u.full_name)
LEFT JOIN clients assigned_client ON np.assign_to = assigned_client.id
WHERE u.id = sqlc.arg("id")
GROUP BY u.id, u.full_name, u.role, b.name;

-- name: GetClientAdminsReportData :many
SELECT 
    c.full_name AS name,
    b.name AS branch_name,
    COUNT(DISTINCT CASE WHEN l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) AS total_loan_given,
    COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) AS defaulted_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'ACTIVE' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) AS active_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'COMPLETED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) AS completed_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'INACTIVE' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) AS inactive_loans,
    COALESCE(c.overpayment, 0) AS overpayment,
    c.phone_number,
    COALESCE(
        (SELECT SUM(np.amount) FROM non_posted np 
         WHERE np.assign_to = c.id AND np.paid_date BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")), 
        0
    ) AS total_paid,
    COALESCE(
        (SELECT SUM(p.loan_amount) 
         FROM loans l 
         JOIN products p ON l.product_id = p.id
         WHERE l.client_id = c.id AND l.status != 'INACTIVE'
         AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")), 
        0
    ) AS total_disbursed,
    COALESCE(
        (SELECT SUM(p.repay_amount - l.paid_amount) 
         FROM loans l 
         JOIN products p ON l.product_id = p.id
         WHERE l.client_id = c.id AND l.status IN ('ACTIVE', 'DEFAULTED')
         AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")), 
        0
    ) AS total_owed,
    COALESCE(
        ((COUNT(DISTINCT CASE WHEN l.status = 'COMPLETED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) - 
         COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END)) * 100.0)
         / NULLIF(COUNT(DISTINCT CASE WHEN l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END), 0), 
        0
    ) AS rate_score,
    COALESCE(
        (COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END) * 100.0) 
        / NULLIF(COUNT(DISTINCT CASE WHEN l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date") THEN l.id END), 0), 
        0
    ) AS default_rate

FROM clients c
LEFT JOIN branches b ON c.branch_id = b.id
LEFT JOIN loans l ON l.client_id = c.id
LEFT JOIN products p ON l.product_id = p.id

GROUP BY c.id, c.full_name, b.name, c.phone_number, c.overpayment
ORDER BY c.full_name;

-- name: GetClientClientsReportData :one
SELECT 
    c.full_name AS name,
    c.phone_number,
    c.id_number,
    c.dob,
    b.name AS branch_name,
    u.full_name AS assigned_staff,
    c.active,
    (
        SELECT COALESCE(
            JSON_ARRAYAGG(
                JSON_OBJECT(
                    'loan_id', l.id,
                    'status', l.status,
                    'loan_amount', p.loan_amount,
                    'repay_amount', p.repay_amount,
                    'paid_amount', l.paid_amount,
                    'disbursed_on', l.disbursed_on,
                    'transaction_fee', l.fee_paid,
                    'approved_by', approved_by_user.full_name
                )
            ), '[]'
        )
        FROM loans l
        JOIN products p ON l.product_id = p.id
        LEFT JOIN users approved_by_user ON l.approved_by = approved_by_user.id
        WHERE l.client_id = c.id
        AND l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")
    ) AS loans,
    (
        SELECT COALESCE(
            JSON_ARRAYAGG(
                JSON_OBJECT(
                    'transaction_number', np.transaction_number,
                    'transaction_source', np.transaction_source,
                    'account_number', np.account_number,
                    'paying_name', np.paying_name,
                    'assigned_by', np.assigned_by,
                    'amount_paid', np.amount,
                    'paid_date', np.paid_date
                )
            ), '[]'
        )
        FROM non_posted np
        LEFT JOIN users assigned_by_user ON np.assigned_by = assigned_by_user.id
        WHERE np.assign_to = c.id
        AND np.paid_date BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")
    ) AS payments

FROM clients c
LEFT JOIN branches b ON c.branch_id = b.id
LEFT JOIN users u ON c.assigned_staff = u.id

WHERE c.id = sqlc.arg("id")
GROUP BY c.id, c.full_name, c.phone_number, c.id_number, c.dob, b.name, u.full_name, c.active;

-- name: GetProductReportData :many
SELECT 
    CONCAT(p.loan_amount, ' - ', b.name) AS product_name,
    COUNT(l.id) AS total_loans_issued,
    COALESCE(SUM(CASE WHEN l.status != 'INACTIVE' THEN p.loan_amount ELSE 0 END), 0) AS total_amount_disbursed,
    COALESCE(SUM(l.paid_amount), 0) AS total_amount_repaid,
    COALESCE(SUM(CASE WHEN l.status != 'INACTIVE' THEN (p.repay_amount - l.paid_amount) ELSE 0 END), 0) AS total_outstanding_amount,
    COUNT(CASE WHEN l.status = 'ACTIVE' THEN l.id END) AS active_loans,
    COUNT(CASE WHEN l.status = 'COMPLETED' THEN l.id END) AS completed_loans,
    COUNT(CASE WHEN l.status = 'DEFAULTED' THEN l.id END) AS defaulted_loans,
    COALESCE(
        (COUNT(CASE WHEN l.status = 'DEFAULTED' THEN 1 END) * 100.0) / NULLIF(COUNT(CASE WHEN l.status != 'INACTIVE' THEN l.id END), 0),
        0
    ) AS default_rate

FROM products p
LEFT JOIN loans l ON l.product_id = p.id 
    AND l.status != 'INACTIVE'  
    AND l.disbursed_on BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")  
LEFT JOIN branches b ON p.branch_id = b.id 

GROUP BY p.loan_amount, b.name
ORDER BY total_loans_issued DESC;

-- name: GetLoansReportData :many
SELECT 
    l.id AS loan_id,
    c.full_name AS client_name,
    b.name AS branch_name,
    u.full_name AS loan_officer,
    p.loan_amount,
    p.repay_amount,
    l.paid_amount,
    COALESCE(p.repay_amount - l.paid_amount, 0) AS outstanding_amount,
    l.status,
    l.total_installments,
    COUNT(CASE WHEN i.paid = TRUE THEN i.id END) AS paid_installments,
    l.due_date AS due_date,
    l.disbursed_on AS disbursed_date,
    COALESCE(
        (COUNT(CASE WHEN i.paid = TRUE THEN 1 END) * 100.0) / NULLIF(l.total_installments, 0),
        0
    ) AS default_risk

FROM loans l
JOIN clients c ON l.client_id = c.id
JOIN branches b ON c.branch_id = b.id
JOIN users u ON l.loan_officer = u.id
JOIN products p ON l.product_id = p.id
LEFT JOIN installments i ON i.loan_id = l.id

WHERE l.created_at BETWEEN sqlc.arg("start_date") AND sqlc.arg("end_date")

GROUP BY l.id, c.full_name, b.name, u.full_name, p.loan_amount, p.repay_amount, l.paid_amount, l.status, l.total_installments
ORDER BY l.created_at DESC;

-- name: GetLoanReportDataById :one
SELECT 
    l.id AS loan_id,
    c.full_name AS client_name,
    p.loan_amount,
    p.repay_amount,
    l.paid_amount,
    l.status,
    l.total_installments,
    COUNT(CASE WHEN i.paid = TRUE THEN i.id END) AS paid_installments,
    (l.total_installments - COUNT(CASE WHEN i.paid = TRUE THEN i.id END)) AS remaining_installments,
    COALESCE(
        JSON_ARRAYAGG(
            JSON_OBJECT(
                'installment_number', i.installment_number,
                'amount_due', i.amount_due,
                'remaining_amount', i.remaining_amount,
                'due_date', i.due_date,
                'paid', i.paid,
                'paid_at', i.paid_at
            )
        ), '[]'
    ) AS installment_details

FROM loans l
JOIN clients c ON l.client_id = c.id
JOIN products p ON l.product_id = p.id
LEFT JOIN installments i ON i.loan_id = l.id

WHERE l.id = sqlc.arg("id")

GROUP BY l.id, c.full_name, p.loan_amount, p.repay_amount, l.paid_amount, l.status, l.total_installments;
