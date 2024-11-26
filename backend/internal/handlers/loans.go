package handlers

import (
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type loanResponse struct {
	ID                 uint32
	Product            loanProductResponse
	Client             loanClientResponse
	LoanOfficer        loanUserResponse
	LoanPurpose        string
	DueDate            time.Time
	ApprovedBy         loanUserResponse
	DisbursedOn        time.Time
	DisbursedBy        loanUserResponse
	TotalInstallments  uint32
	InstallmentsPeriod uint32
	Status             string
	ProcessingFee      float64
	PaidAmount         float64
	UpdatedBy          loanUserResponse
	CreatedBy          loanUserResponse
	CreatedAt          time.Time
}

type loanProductResponse struct {
	ID             uint32  `json:"id"`
	BranchName     string  `json:"branch_name"`
	LoanAmount     float64 `json:"loan_amount"`
	RepayAmount    float64 `json:"repay_amount"`
	InterestAmount float64 `json:"interest_amount"`
}

type loanClientResponse struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	IdNumber    string    `json:"id_number"`
	Dob         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Active      bool      `json:"active"`
	BranchName  string    `json:"branch_name"`
	Overpayment float64   `json:"overpayment"`
}

type loanUserResponse struct {
	ID       uint32 `json:"id"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
}

func (s *Server) createLoan(ctx *gin.Context) {}

func (s *Server) disburseLoan(ctx *gin.Context) {}

func (s *Server) transferLoanOfficer(ctx *gin.Context) {}

func (s *Server) listLoans(ctx *gin.Context) {}

func (s *Server) getLoan(ctx *gin.Context) {}

func (s *Server) listLoansByClient(ctx *gin.Context) {}

func (s *Server) listLoansByBranch(ctx *gin.Context) {}

func (s *Server) listLoansByLoanOfficer(ctx *gin.Context) {}

func (s *Server) listLoansByStatus(ctx *gin.Context) {}

func (s *Server) listNonDisbursedLoans(ctx *gin.Context) {}

func (s *Server) updateLoanStatusByAdmin(ctx *gin.Context) {}

func (s *Server) structureLoan(loan *repository.Loan, ctx *gin.Context) (loanResponse, error) {
	product, err := s.repo.Products.GetProductByID(ctx, loan.ProductID)
	if err != nil {
		return loanResponse{}, err
	}

	productBranch, err := s.repo.Branches.GetBranchByID(ctx, product.BranchID)
	if err != nil {
		return loanResponse{}, err
	}

	client, err := s.repo.Clients.GetClient(ctx, loan.ClientID)
	if err != nil {
		return loanResponse{}, err
	}

	clientBranch, err := s.repo.Branches.GetBranchByID(ctx, client.BranchID)
	if err != nil {
		return loanResponse{}, err
	}

	loanOfficer, err := s.repo.Users.GetUserByID(ctx, loan.LoanOfficerID)
	if err != nil {
		return loanResponse{}, err
	}

	approvedBy, err := s.repo.Users.GetUserByID(ctx, loan.ApprovedBy)
	if err != nil {
		return loanResponse{}, err
	}

	disbursedBy := loanUserResponse{}

	if loan.DisbursedBy != nil {
		disbursedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.DisbursedBy)
		if err != nil {
			return loanResponse{}, err
		}

		disbursedBy = loanUserResponse{
			ID:       disbursedByUser.ID,
			Fullname: disbursedByUser.FullName,
			Role:     disbursedByUser.Role,
		}
	}

	updatedBy := loanUserResponse{}

	if loan.UpdatedBy != nil {
		updatedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.UpdatedBy)
		if err != nil {
			return loanResponse{}, err
		}

		updatedBy = loanUserResponse{
			ID:       updatedByUser.ID,
			Fullname: updatedByUser.FullName,
			Role:     updatedByUser.Role,
		}
	}

	createdBy, err := s.repo.Users.GetUserByID(ctx, loan.CreatedBy)
	if err != nil {
		return loanResponse{}, err
	}

	return loanResponse{
		ID: loan.ID,
		Product: loanProductResponse{
			ID:             product.ID,
			BranchName:     productBranch.Name,
			LoanAmount:     product.LoanAmount,
			RepayAmount:    product.RepayAmount,
			InterestAmount: product.InterestAmount,
		},
		Client: loanClientResponse{
			ID:          client.ID,
			Name:        client.FullName,
			PhoneNumber: client.PhoneNumber,
			IdNumber:    *client.IdNumber,
			Dob:         *client.Dob,
			Gender:      client.Gender,
			Active:      client.Active,
			BranchName:  clientBranch.Name,
			Overpayment: client.Overpayment,
		},
		LoanOfficer:        loanUserResponse{ID: loanOfficer.ID, Fullname: loanOfficer.FullName, Role: loanOfficer.Role},
		LoanPurpose:        *loan.LoanPurpose,
		DueDate:            *loan.DueDate,
		ApprovedBy:         loanUserResponse{ID: approvedBy.ID, Fullname: approvedBy.FullName, Role: approvedBy.Role},
		DisbursedOn:        *loan.DisbursedOn,
		DisbursedBy:        disbursedBy,
		TotalInstallments:  loan.TotalInstallments,
		InstallmentsPeriod: loan.InstallmentsPeriod,
		Status:             loan.Status,
		ProcessingFee:      loan.ProcessingFee,
		PaidAmount:         loan.PaidAmount,
		UpdatedBy:          updatedBy,
		CreatedBy:          loanUserResponse{ID: createdBy.ID, Fullname: createdBy.FullName, Role: createdBy.Role},
		CreatedAt:          loan.CreatedAt,
	}, nil
}
