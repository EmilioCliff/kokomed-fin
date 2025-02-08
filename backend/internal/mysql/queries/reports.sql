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