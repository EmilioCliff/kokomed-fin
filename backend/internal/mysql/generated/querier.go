// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package generated

import (
	"context"
	"database/sql"
)

type Querier interface {
	AssignNonPosted(ctx context.Context, arg AssignNonPostedParams) (sql.Result, error)
	CheckActiveLoanForClient(ctx context.Context, clientID uint32) (bool, error)
	CreateBranch(ctx context.Context, name string) (sql.Result, error)
	CreateClient(ctx context.Context, arg CreateClientParams) (sql.Result, error)
	CreateInstallment(ctx context.Context, arg CreateInstallmentParams) (sql.Result, error)
	CreateLoan(ctx context.Context, arg CreateLoanParams) (sql.Result, error)
	CreateNonPosted(ctx context.Context, arg CreateNonPostedParams) (sql.Result, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	DeleteBranch(ctx context.Context, id uint32) error
	DeleteClient(ctx context.Context, id uint32) (sql.Result, error)
	DeleteLoan(ctx context.Context, id uint32) error
	DeleteNonPosted(ctx context.Context, id uint32) error
	DeleteProduct(ctx context.Context, id uint32) error
	DisburseLoan(ctx context.Context, arg DisburseLoanParams) (sql.Result, error)
	GetBranch(ctx context.Context, id uint32) (Branch, error)
	GetClient(ctx context.Context, id uint32) (Client, error)
	GetClientActiveLoan(ctx context.Context, arg GetClientActiveLoanParams) (uint32, error)
	GetClientIDByPhoneNumber(ctx context.Context, phoneNumber string) (uint32, error)
	GetInstallment(ctx context.Context, id uint32) (Installment, error)
	GetLoan(ctx context.Context, id uint32) (Loan, error)
	GetLoanPaymentData(ctx context.Context, id uint32) (GetLoanPaymentDataRow, error)
	GetNonPosted(ctx context.Context, id uint32) (NonPosted, error)
	GetProduct(ctx context.Context, id uint32) (Product, error)
	GetProductRepayAmount(ctx context.Context, id uint32) (float64, error)
	GetUser(ctx context.Context, id uint32) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListAllNonPosted(ctx context.Context, arg ListAllNonPostedParams) ([]NonPosted, error)
	ListAllNonPostedByTransactionSource(ctx context.Context, transactionSource NonPostedTransactionSource) ([]NonPosted, error)
	ListBranches(ctx context.Context) ([]Branch, error)
	ListClients(ctx context.Context, arg ListClientsParams) ([]Client, error)
	ListClientsByActiveStatus(ctx context.Context, arg ListClientsByActiveStatusParams) ([]Client, error)
	ListClientsByBranch(ctx context.Context, arg ListClientsByBranchParams) ([]Client, error)
	ListInstallmentsByLoan(ctx context.Context, arg ListInstallmentsByLoanParams) ([]Installment, error)
	ListLoans(ctx context.Context, arg ListLoansParams) ([]ListLoansRow, error)
	ListLoansByClient(ctx context.Context, arg ListLoansByClientParams) ([]Loan, error)
	ListLoansByLoanOfficer(ctx context.Context, arg ListLoansByLoanOfficerParams) ([]Loan, error)
	ListLoansByStatus(ctx context.Context, arg ListLoansByStatusParams) ([]Loan, error)
	ListNonDisbursedLoans(ctx context.Context, arg ListNonDisbursedLoansParams) ([]Loan, error)
	ListNonPostedByTransactionSource(ctx context.Context, arg ListNonPostedByTransactionSourceParams) ([]NonPosted, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListProductsByBranch(ctx context.Context, arg ListProductsByBranchParams) ([]Product, error)
	ListUnassignedNonPosted(ctx context.Context, arg ListUnassignedNonPostedParams) ([]NonPosted, error)
	ListUnpaidInstallmentsByLoan(ctx context.Context, loanID uint32) ([]Installment, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	PayInstallment(ctx context.Context, arg PayInstallmentParams) (sql.Result, error)
	TransferLoan(ctx context.Context, arg TransferLoanParams) (sql.Result, error)
	UpdateBranch(ctx context.Context, arg UpdateBranchParams) (sql.Result, error)
	UpdateClient(ctx context.Context, arg UpdateClientParams) (sql.Result, error)
	// -- name: UpdateClientOverpayment :execresult
	// UPDATE clients
	//     SET overpayment = overpayment + sqlc.arg("overpayment")
	// WHERE phone_number = sqlc.arg("phone_number");
	UpdateClientOverpayment(ctx context.Context, arg UpdateClientOverpaymentParams) (sql.Result, error)
	UpdateInstallment(ctx context.Context, arg UpdateInstallmentParams) (sql.Result, error)
	UpdateLoan(ctx context.Context, arg UpdateLoanParams) (sql.Result, error)
	UpdateLoanProcessingFeeStatus(ctx context.Context, arg UpdateLoanProcessingFeeStatusParams) (sql.Result, error)
	UpdateLoanStatus(ctx context.Context, arg UpdateLoanStatusParams) (sql.Result, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
