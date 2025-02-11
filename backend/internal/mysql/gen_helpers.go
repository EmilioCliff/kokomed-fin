package mysql

import (
	"context"
	"database/sql"
	"time"
)

const getUserAdminsReportData = `-- name: GetUserAdminsReportData :many
SELECT 
    u.full_name AS name,
    u.role,
    b.name AS branch_name,
    COUNT(DISTINCT CASE 
        WHEN l.approved_by = u.id AND l.created_at BETWEEN ? AND ? THEN l.id 
    END) AS approved_loans,
    COUNT(DISTINCT CASE 
        WHEN l.loan_officer = u.id AND l.status = 'ACTIVE' AND l.created_at BETWEEN ? AND ? THEN l.id 
    END) AS active_loans,
    COUNT(DISTINCT CASE 
        WHEN l.loan_officer = u.id AND l.status = 'COMPLETED' AND l.created_at BETWEEN ? AND ? THEN l.id 
    END) AS completed_loans,
    COALESCE(
        (COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id AND l.status = 'DEFAULTED' AND l.created_at BETWEEN ? AND ? THEN l.id 
        END) * 100.0) / NULLIF(COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id AND l.created_at BETWEEN ? AND ? THEN l.id 
        END), 0), 
        0
    ) AS default_rate,
    COUNT(DISTINCT CASE 
        WHEN c.created_by = u.id AND c.created_at BETWEEN ? AND ? THEN c.id 
    END) AS clients_registered,
    COUNT(DISTINCT CASE 
        WHEN LOWER(np.assigned_by) = LOWER(u.full_name) AND np.paid_date BETWEEN ? AND ? THEN np.id 
    END) AS payments_assigned
FROM users u
LEFT JOIN branches b ON u.branch_id = b.id
LEFT JOIN loans l ON l.loan_officer = u.id
LEFT JOIN clients c ON c.created_by = u.id
LEFT JOIN non_posted np ON LOWER(np.assigned_by) = LOWER(u.full_name)
GROUP BY u.id, u.full_name, u.role, b.name
ORDER BY u.role DESC, u.full_name
`

type GetUserAdminsReportDataParams struct {
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	ActiveStart    time.Time `json:"active_start"`
	ActiveEnd      time.Time `json:"active_end"`
	CompletedStart time.Time `json:"completed_start"`
	CompletedEnd   time.Time `json:"completed_end"`
	DefaultedStart time.Time `json:"defaulted_start"`
	DefaultedEnd   time.Time `json:"defaulted_end"`
	TotalStart     time.Time `json:"total_start"`
	TotalEnd       time.Time `json:"total_end"`
	ClientsStart   time.Time `json:"clients_start"`
	ClientsEnd     time.Time `json:"clients_end"`
	PaymentsStart  time.Time `json:"payments_start"`
	PaymentsEnd    time.Time `json:"payments_end"`
}

type GetUserAdminsReportDataRow struct {
	Name              string         `json:"name"`
	Role              string      `json:"role"`
	BranchName        sql.NullString `json:"branch_name"`
	ApprovedLoans     int64          `json:"approved_loans"`
	ActiveLoans       int64          `json:"active_loans"`
	CompletedLoans    int64          `json:"completed_loans"`
	DefaultRate       interface{}    `json:"default_rate"`
	ClientsRegistered int64          `json:"clients_registered"`
	PaymentsAssigned  int64          `json:"payments_assigned"`
}

func (q *UserRepository) GetUserAdminsReportData(ctx context.Context, arg GetUserAdminsReportDataParams) ([]GetUserAdminsReportDataRow, error) {
	rows, err := q.db.db.QueryContext(ctx, getUserAdminsReportData,
		arg.StartDate,
		arg.EndDate,
		arg.ActiveStart,
		arg.ActiveEnd,
		arg.CompletedStart,
		arg.CompletedEnd,
		arg.DefaultedStart,
		arg.DefaultedEnd,
		arg.TotalStart,
		arg.TotalEnd,
		arg.ClientsStart,
		arg.ClientsEnd,
		arg.PaymentsStart,
		arg.PaymentsEnd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserAdminsReportDataRow{}
	for rows.Next() {
		var i GetUserAdminsReportDataRow
		if err := rows.Scan(
			&i.Name,
			&i.Role,
			&i.BranchName,
			&i.ApprovedLoans,
			&i.ActiveLoans,
			&i.CompletedLoans,
			&i.DefaultRate,
			&i.ClientsRegistered,
			&i.PaymentsAssigned,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserUsersReportData = `-- name: GetUserUsersReportData :many
SELECT 
    u.full_name AS name,
    u.role,
    b.name AS branch,
    COUNT(DISTINCT CASE 
        WHEN c.assigned_staff = u.id AND c.created_at BETWEEN ? AND ? THEN c.id 
    END) AS total_clients_handled,
    COUNT(DISTINCT CASE 
        WHEN l.loan_officer = u.id AND l.created_at BETWEEN ? AND ? THEN l.id 
    END) AS loans_approved,
    COALESCE(SUM(CASE 
        WHEN l.loan_officer = u.id AND l.created_at BETWEEN ? AND ? THEN p.repay_amount 
    END), 0) AS total_loan_amount_managed,
    COALESCE(SUM(CASE 
        WHEN l.loan_officer = u.id AND l.paid_amount > 0 AND l.created_at BETWEEN ? AND ? THEN l.paid_amount 
    END), 0) AS total_collected_amount,
    COALESCE(
        (COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id AND l.status = 'DEFAULTED' AND l.created_at BETWEEN ? AND ? THEN l.id 
        END) * 100.0) / NULLIF(COUNT(DISTINCT CASE 
            WHEN l.loan_officer = u.id AND l.created_at BETWEEN ? AND ? THEN l.id 
        END), 0), 
        0
    ) AS default_rate,
    COUNT(DISTINCT CASE 
        WHEN LOWER(np.assigned_by) = LOWER(u.full_name) AND np.paid_date BETWEEN ? AND ? THEN np.id 
    END) AS assigned_payments

FROM users u
LEFT JOIN branches b ON u.branch_id = b.id
LEFT JOIN clients c ON c.assigned_staff = u.id
LEFT JOIN loans l ON l.loan_officer = u.id
LEFT JOIN products p ON l.product_id = p.id
LEFT JOIN non_posted np ON LOWER(np.assigned_by) = LOWER(u.full_name)

WHERE u.id = ?
GROUP BY u.id, u.full_name, b.name
`

type GetUserUsersReportDataParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ID        uint32    `json:"id"`
}

type GetUserUsersReportDataRow struct {
	Name                   string         `json:"name"`
	Role                   string      `json:"role"`
	Branch                 sql.NullString `json:"branch"`
	TotalClientsHandled    int64          `json:"total_clients_handled"`
	LoansApproved          int64          `json:"loans_approved"`
	TotalLoanAmountManaged interface{}    `json:"total_loan_amount_managed"`
	TotalCollectedAmount   interface{}    `json:"total_collected_amount"`
	DefaultRate            interface{}    `json:"default_rate"`
	AssignedPayments       int64          `json:"assigned_payments"`
}

func (q *UserRepository) GetUserUsersReportData(ctx context.Context, arg GetUserUsersReportDataParams) ([]GetUserUsersReportDataRow, error) {
	rows, err := q.db.db.QueryContext(ctx, getUserUsersReportData,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserUsersReportDataRow{}
	for rows.Next() {
		var i GetUserUsersReportDataRow
		if err := rows.Scan(
			&i.Name,
			&i.Role,
			&i.Branch,
			&i.TotalClientsHandled,
			&i.LoansApproved,
			&i.TotalLoanAmountManaged,
			&i.TotalCollectedAmount,
			&i.DefaultRate,
			&i.AssignedPayments,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getClientAdminsReportData = `-- name: GetClientAdminsReportData :many
SELECT 
    c.full_name AS name,
    b.name AS branch_name,
    COUNT(DISTINCT CASE WHEN l.created_at BETWEEN ? AND ? THEN l.id END) AS total_loan_given,
    COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN ? AND ? THEN l.id END) AS defaulted_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'ACTIVE' AND l.created_at BETWEEN ? AND ? THEN l.id END) AS active_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'COMPLETED' AND l.created_at BETWEEN ? AND ? THEN l.id END) AS completed_loans,
    COUNT(DISTINCT CASE WHEN l.status = 'INACTIVE' AND l.created_at BETWEEN ? AND ? THEN l.id END) AS inactive_loans,
    COALESCE(c.overpayment, 0) AS overpayment,
    c.phone_number,
    COALESCE(
        (SELECT SUM(np.amount) FROM non_posted np 
         WHERE np.assign_to = c.id AND np.paid_date BETWEEN ? AND ?), 
        0
    ) AS total_paid,
    COALESCE(
        (SELECT SUM(p.loan_amount) 
         FROM loans l 
         JOIN products p ON l.product_id = p.id
         WHERE l.client_id = c.id AND l.status != 'INACTIVE'
         AND l.created_at BETWEEN ? AND ?), 
        0
    ) AS total_disbursed,
    COALESCE(
        (SELECT SUM(p.repay_amount - l.paid_amount) 
         FROM loans l 
         JOIN products p ON l.product_id = p.id
         WHERE l.client_id = c.id AND l.status IN ('ACTIVE', 'DEFAULTED')
         AND l.created_at BETWEEN ? AND ?), 
        0
    ) AS total_owed,
    COALESCE(
        ((COUNT(DISTINCT CASE WHEN l.status = 'COMPLETED' AND l.created_at BETWEEN ? AND ? THEN l.id END) - 
         COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN ? AND ? THEN l.id END)) * 100.0)
         / NULLIF(COUNT(DISTINCT CASE WHEN l.created_at BETWEEN ? AND ? THEN l.id END), 0), 
        0
    ) AS rate_score,
    COALESCE(
        (COUNT(DISTINCT CASE WHEN l.status = 'DEFAULTED' AND l.created_at BETWEEN ? AND ? THEN l.id END) * 100.0) 
        / NULLIF(COUNT(DISTINCT CASE WHEN l.created_at BETWEEN ? AND ? THEN l.id END), 0), 
        0
    ) AS default_rate

FROM clients c
LEFT JOIN branches b ON c.branch_id = b.id
LEFT JOIN loans l ON l.client_id = c.id
LEFT JOIN products p ON l.product_id = p.id

GROUP BY c.id, c.full_name, b.name, c.phone_number, c.overpayment
ORDER BY c.full_name
`

type GetClientAdminsReportDataParams struct {
	EndDate   time.Time `json:"end_date"`
	StartDate time.Time `json:"start_date"`
}

type GetClientAdminsReportDataRow struct {
	Name           string         `json:"name"`
	BranchName     sql.NullString `json:"branch_name"`
	TotalLoanGiven int64          `json:"total_loan_given"`
	DefaultedLoans int64          `json:"defaulted_loans"`
	ActiveLoans    int64          `json:"active_loans"`
	CompletedLoans int64          `json:"completed_loans"`
	InactiveLoans  int64          `json:"inactive_loans"`
	Overpayment    float64        `json:"overpayment"`
	PhoneNumber    string         `json:"phone_number"`
	TotalPaid      interface{}    `json:"total_paid"`
	TotalDisbursed interface{}    `json:"total_disbursed"`
	TotalOwed      interface{}    `json:"total_owed"`
	RateScore      interface{}    `json:"rate_score"`
	DefaultRate    interface{}    `json:"default_rate"`
}

func (q *ClientRepository) GetClientAdminsReportData(ctx context.Context, arg GetClientAdminsReportDataParams) ([]GetClientAdminsReportDataRow, error) {
	rows, err := q.db.db.QueryContext(ctx, getClientAdminsReportData,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetClientAdminsReportDataRow{}
	for rows.Next() {
		var i GetClientAdminsReportDataRow
		if err := rows.Scan(
			&i.Name,
			&i.BranchName,
			&i.TotalLoanGiven,
			&i.DefaultedLoans,
			&i.ActiveLoans,
			&i.CompletedLoans,
			&i.InactiveLoans,
			&i.Overpayment,
			&i.PhoneNumber,
			&i.TotalPaid,
			&i.TotalDisbursed,
			&i.TotalOwed,
			&i.RateScore,
			&i.DefaultRate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getClientClientsReportData = `-- name: GetClientClientsReportData :one
SELECT 
    c.full_name AS name,
    c.phone_number,
    c.id_number,
    c.dob,
    b.name AS branch_name,
    u.full_name AS assigned_staff,
    c.active,

    -- ✅ Loan Details (Get all loans separately to prevent duplication)
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
                    'created_by', created_by_user.full_name,
                    'assigned_by', approved_by_user.full_name
                )
            ), '[]'
        )
        FROM loans l
        JOIN products p ON l.product_id = p.id
        LEFT JOIN users created_by_user ON l.created_by = created_by_user.id
        LEFT JOIN users approved_by_user ON l.approved_by = approved_by_user.id
        WHERE l.client_id = c.id
        AND l.created_at BETWEEN ? AND ?
    ) AS loans,

    -- ✅ Payment Details (Get all payments separately to prevent duplication)
    (
        SELECT COALESCE(
            JSON_ARRAYAGG(
                JSON_OBJECT(
                    'transaction_number', np.transaction_number,
                    'transaction_source', np.transaction_source,
                    'account_number', np.account_number,
                    'paying_name', np.paying_name,
                    'assigned_by', assigned_by_user.full_name,
                    'amount_paid', np.amount,
                    'paid_date', np.paid_date
                )
            ), '[]'
        )
        FROM non_posted np
        LEFT JOIN users assigned_by_user ON np.assigned_by = assigned_by_user.id
        WHERE np.assign_to = c.id
        AND np.paid_date BETWEEN ? AND ?
    ) AS payments

FROM clients c
LEFT JOIN branches b ON c.branch_id = b.id
LEFT JOIN users u ON c.assigned_staff = u.id

WHERE c.id = ?
GROUP BY c.id, c.full_name, c.phone_number, c.id_number, c.dob, b.name, u.full_name, c.active
`

type GetClientClientsReportDataParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ID        uint32    `json:"id"`
}

type GetClientClientsReportDataRow struct {
	Name          string         `json:"name"`
	PhoneNumber   string         `json:"phone_number"`
	IDNumber      sql.NullString `json:"id_number"`
	Dob           sql.NullTime   `json:"dob"`
	BranchName    sql.NullString `json:"branch_name"`
	AssignedStaff sql.NullString `json:"assigned_staff"`
	Active        bool           `json:"active"`
	Loans         interface{}    `json:"loans"`
	Payments      interface{}    `json:"payments"`
}


func (q *ClientRepository) GetClientClientsReportData(ctx context.Context, arg GetClientClientsReportDataParams) (GetClientClientsReportDataRow, error) {
	row := q.db.db.QueryRowContext(ctx, getClientClientsReportData,
		arg.StartDate,
		arg.EndDate,
		arg.StartDate,
		arg.EndDate,
		arg.ID,
	)
	var i GetClientClientsReportDataRow
	err := row.Scan(
		&i.Name,
		&i.PhoneNumber,
		&i.IDNumber,
		&i.Dob,
		&i.BranchName,
		&i.AssignedStaff,
		&i.Active,
		&i.Loans,
		&i.Payments,
	)
	return i, err
}

const getProductReportData = `-- name: GetProductReportData :many
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
    AND l.created_at BETWEEN ? AND ?  
LEFT JOIN branches b ON p.branch_id = b.id 

GROUP BY p.loan_amount, b.name
ORDER BY total_loans_issued DESC
`

type GetProductReportDataParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type GetProductReportDataRow struct {
	ProductName            string      `json:"product_name"`
	TotalLoansIssued       int64       `json:"total_loans_issued"`
	TotalAmountDisbursed   interface{} `json:"total_amount_disbursed"`
	TotalAmountRepaid      interface{} `json:"total_amount_repaid"`
	TotalOutstandingAmount interface{} `json:"total_outstanding_amount"`
	ActiveLoans            int64       `json:"active_loans"`
	CompletedLoans         int64       `json:"completed_loans"`
	DefaultedLoans         int64       `json:"defaulted_loans"`
	DefaultRate            interface{} `json:"default_rate"`
}

func (q *ProductRepository) GetProductReportData(ctx context.Context, arg GetProductReportDataParams) ([]GetProductReportDataRow, error) {
	rows, err := q.db.db.QueryContext(ctx, getProductReportData,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetProductReportDataRow{}
	for rows.Next() {
		var i GetProductReportDataRow
		if err := rows.Scan(
			&i.ProductName,
			&i.TotalLoansIssued,
			&i.TotalAmountDisbursed,
			&i.TotalAmountRepaid,
			&i.TotalOutstandingAmount,
			&i.ActiveLoans,
			&i.CompletedLoans,
			&i.DefaultedLoans,
			&i.DefaultRate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLoansReportData = `-- name: GetLoansReportData :many
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

WHERE l.created_at BETWEEN ? AND ?

GROUP BY l.id, c.full_name, b.name, u.full_name, p.loan_amount, p.repay_amount, l.paid_amount, l.status, l.total_installments
ORDER BY l.created_at DESC
`

type GetLoansReportDataParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type GetLoansReportDataRow struct {
	LoanID            uint32      `json:"loan_id"`
	ClientName        string      `json:"client_name"`
	BranchName        string      `json:"branch_name"`
	LoanOfficer       string      `json:"loan_officer"`
	LoanAmount        float64     `json:"loan_amount"`
	RepayAmount       float64     `json:"repay_amount"`
	PaidAmount        float64     `json:"paid_amount"`
	OutstandingAmount interface{} `json:"outstanding_amount"`
	Status            string 	  `json:"status"`
	TotalInstallments uint32      `json:"total_installments"`
	PaidInstallments  int64       `json:"paid_installments"`
	DueDate           sql.NullTime `json:"due_date"`
	DisbursedDate 	  sql.NullTime `json:"disbursed_date"`
	DefaultRisk       interface{} `json:"default_risk"`
}

func (q *LoanRepository) GetLoansReportData(ctx context.Context, arg GetLoansReportDataParams) ([]GetLoansReportDataRow, error) {
	rows, err := q.db.db.QueryContext(ctx, getLoansReportData,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLoansReportDataRow{}
	for rows.Next() {
		var i GetLoansReportDataRow
		if err := rows.Scan(
			&i.LoanID,
			&i.ClientName,
			&i.BranchName,
			&i.LoanOfficer,
			&i.LoanAmount,
			&i.RepayAmount,
			&i.PaidAmount,
			&i.OutstandingAmount,
			&i.Status,
			&i.TotalInstallments,
			&i.PaidInstallments,
			&i.DueDate,
			&i.DisbursedDate,
			&i.DefaultRisk,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLoanReportDataById = `-- name: GetLoanReportDataById :one
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
                'installment_amount', i.amount_due,
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

WHERE l.id = ?

GROUP BY l.id, c.full_name, p.loan_amount, p.repay_amount, l.paid_amount, l.status, l.total_installments
`

type GetLoanReportDataByIdRow struct {
	LoanID                uint32      `json:"loan_id"`
	ClientName            string      `json:"client_name"`
	LoanAmount            float64     `json:"loan_amount"`
	RepayAmount           float64     `json:"repay_amount"`
	PaidAmount            float64     `json:"paid_amount"`
	Status                string 	  `json:"status"`
	TotalInstallments     uint32      `json:"total_installments"`
	PaidInstallments      int64       `json:"paid_installments"`
	RemainingInstallments uint32      `json:"remaining_installments"`
	InstallmentDetails    interface{} `json:"installment_details"`
}

func (q *LoanRepository) GetLoanReportDataById(ctx context.Context, id uint32) (GetLoanReportDataByIdRow, error) {
	row := q.db.db.QueryRowContext(ctx, getLoanReportDataById, id)
	var i GetLoanReportDataByIdRow
	err := row.Scan(
		&i.LoanID,
		&i.ClientName,
		&i.LoanAmount,
		&i.RepayAmount,
		&i.PaidAmount,
		&i.Status,
		&i.TotalInstallments,
		&i.PaidInstallments,
		&i.RemainingInstallments,
		&i.InstallmentDetails,
	)
	return i, err
}
