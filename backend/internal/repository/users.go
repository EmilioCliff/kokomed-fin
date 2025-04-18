package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type User struct {
	ID              uint32    `json:"id"`
	FullName        string    `json:"fullName"`
	PhoneNumber     string    `json:"phoneNumber"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PasswordUpdated uint32    `json:"password_updated"`
	RefreshToken    string    `json:"refresh_token"`
	Role            string    `json:"role"`
	BranchName      *string   `json:"branchName"`
	BranchID        uint32    `json:"branch_id"`
	UpdatedBy       uint32    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uint32    `json:"created_by"`
	CreatedAt       time.Time `json:"createdAt"`
}

type UpdateUser struct {
	ID           uint32     `json:"id"`
	Role         *string    `json:"role"`
	BranchID     *uint32    `json:"branch_id"`
	Password     *string    `json:"password"`
	RefreshToken *string    `json:"refresh_token"`
	UpdatedBy    *uint32    `json:"updated_by"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type UserShortResponse struct {
	ID          uint32 `json:"id"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

type CategorySearch struct {
	Search *string `json:"search"`
	Role   *string `json:"role"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (User, error)
	GetUserByID(ctx context.Context, id uint32) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListUsers(
		ctx context.Context,
		category *CategorySearch,
		pgData *pkg.PaginationMetadata,
	) ([]User, pkg.PaginationMetadata, error)
	UpdateUser(ctx context.Context, user *UpdateUser) (User, error)
	UpdateUserPassword(ctx context.Context, email string, password string) error
	CheckUserExistance(ctx context.Context, email string) bool

	GetReportUserAdminData(
		ctx context.Context,
		filters services.ReportFilters,
	) ([]services.UserAdminsReportData, services.UserAdminsSummary, error)
	GetReportUserUsersData(
		ctx context.Context,
		id uint32,
		filters services.ReportFilters,
	) (services.UserUsersReportData, error)
}
