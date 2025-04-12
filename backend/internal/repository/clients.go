package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type Client struct {
	ID            uint32     `json:"id"`
	FullName      string     `json:"full_name"`
	PhoneNumber   string     `json:"phone_number"`
	IdNumber      *string    `json:"id_number"`
	Dob           *time.Time `json:"dob"`
	Gender        string     `json:"gender"`
	Active        bool       `json:"active"`
	BranchID      uint32     `json:"branch_id"`
	AssignedStaff uint32     `json:"assigned_staff"`
	Overpayment   float64    `json:"overpayment"`
	UpdatedBy     uint32     `json:"updated_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     uint32     `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	BranchName    *string    `json:"branch_name"`
	DueAmount     float64    `json:"due_amount"`
}

type UpdateClient struct {
	ID        uint32     `json:"id"`
	UpdatedBy uint32     `json:"updated_by"`
	IdNumber  *string    `json:"id_number"`
	Dob       *time.Time `json:"dob"`
	Active    *bool      `json:"active"`
	BranchID  *uint32    `json:"branch_id"`
}

type ClientCategorySearch struct {
	Search *string `json:"search"`
	Active *bool   `json:"active"`
}

type ClientShort struct {
	ID          uint32  `json:"id"`
	FullName    string  `json:"fullName"`
	PhoneNumber string  `json:"phoneNumber"`
	Active      bool    `json:"active"`
	Overpayment float64 `json:"overpayment"`
	BranchName  string  `json:"branchName"`
	DueAmount   float64 `json:"dueAmount"`
}

type ClientFullData struct {
	ID            uint32            `json:"id"`
	FullName      string            `json:"fullName"`
	PhoneNumber   string            `json:"phoneNumber"`
	IDNumber      *string           `json:"idNumber,omitempty"`
	DOB           *time.Time        `json:"dob,omitempty"`
	Gender        string            `json:"gender"`
	Active        bool              `json:"active"`
	BranchName    string            `json:"branchName"`
	AssignedStaff UserShortResponse `json:"assignedStaff"`
	Overpayment   float64           `json:"overpayment"`
	DueAmount     float64           `json:"dueAmount"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	CreatedBy     UserShortResponse `json:"createdBy"`
	UpdatedBy     UserShortResponse `json:"updatedBy,omitempty"`
}

type ClientRepository interface {
	CreateClient(ctx context.Context, client *Client) (ClientFullData, error)
	UpdateClient(ctx context.Context, client *UpdateClient) error
	UpdateClientOverpayment(ctx context.Context, phoneNumber string, overpayment float64) error
	ListClients(
		ctx context.Context,
		category *ClientCategorySearch,
		pgData *pkg.PaginationMetadata,
	) ([]ClientFullData, pkg.PaginationMetadata, error)
	// GetClient(ctx context.Context, clientID uint32) (ClientFullData, error)
	GetClientFullData(ctx context.Context, clientID uint32) (ClientFullData, error)
	GetClientIDByPhoneNumber(ctx context.Context, phoneNumber string) (uint32, error)
	ListClientsByBranch(
		ctx context.Context,
		branchID uint32,
		pgData *pkg.PaginationMetadata,
	) ([]Client, error)
	ListClientsByActiveStatus(
		ctx context.Context,
		active bool,
		pgData *pkg.PaginationMetadata,
	) ([]Client, error)

	GetReportClientAdminData(
		ctx context.Context,
		filters services.ReportFilters,
	) ([]services.ClientAdminsReportData, services.ClientSummary, error)
	GetReportClientClientsData(
		ctx context.Context,
		id uint32,
		filters services.ReportFilters,
	) (services.ClientClientsReportData, error)
}
