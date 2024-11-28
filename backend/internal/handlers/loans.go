package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
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
	FeePaid            bool
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

type createLoanRequest struct {
	ProductID          uint32  `binding:"required" json:"product_id"`
	ClientID           uint32  `binding:"required" json:"client_id"`
	LoanOfficerID      uint32  `binding:"required" json:"loan_officer_id"`
	LoanPurpose        string  `                   json:"loan_purpose"`
	ApprovedBy         uint32  `binding:"required" json:"approved_by"`
	DisburseBy         uint32  `                   json:"disburse_by"`
	DisburseOn         string  `                   json:"disburse_on"`
	Installments       uint32  `binding:"required" json:"installments"`
	InstallmentsPeriod uint32  `binding:"required" json:"installments_period"`
	ProcessingFee      float64 `binding:"required" json:"processing_fee"`
	ProcessingFeePaid  bool    `                   json:"processing_fee_paid"`
	CreatedBy          uint32  `binding:"required" json:"created_by"`
}

func (s *Server) createLoan(ctx *gin.Context) {
	var req createLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	params := repository.Loan{
		ProductID:          req.ProductID,
		ClientID:           req.ClientID,
		LoanOfficerID:      req.LoanOfficerID,
		ApprovedBy:         req.ApprovedBy,
		TotalInstallments:  req.Installments,
		InstallmentsPeriod: req.InstallmentsPeriod,
		ProcessingFee:      req.ProcessingFee,
		FeePaid:            false,
		CreatedBy:          req.CreatedBy,
	}

	if req.ProcessingFeePaid {
		params.FeePaid = true
	}

	if req.DisburseOn != "" && req.DisburseBy != 0 {
		disburseDate, err := time.Parse("2006-01-02", req.DisburseOn)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid disburse_on date format")))

			return
		}

		params.DisbursedOn = pkg.TimePtr(disburseDate)
		params.DisbursedBy = pkg.Uint32Ptr(req.DisburseBy)

		params.DueDate = pkg.TimePtr(disburseDate.AddDate(0, 0, int(req.Installments)*int(req.InstallmentsPeriod)))
	} else if req.DisburseOn != "" || req.DisburseBy != 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "both disburse_on and disburse_by are required if one is provided")))

		return
	}

	if req.LoanPurpose != "" {
		params.LoanPurpose = pkg.StringPtr(req.LoanPurpose)
	}

	loan, err := s.repo.Loans.CreateLoan(ctx, &params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp, err := s.structureLoan(&loan, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type disburseLoanRequest struct {
	DisburseBy  uint32 `binding:"required" json:"disburse_by"`
	DisbursedOn string `binding:"required" json:"disbursed_on"`
}

func (s *Server) disburseLoan(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var req disburseLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	disbursedDate, err := time.Parse("2006-01-02", req.DisbursedOn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	err = s.repo.Loans.DisburseLoan(ctx, &repository.DisburseLoan{
		DisbursedBy: req.DisburseBy,
		DisbursedOn: disbursedDate,
		ID:          id,
		DueDate:     disbursedDate,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan disbursed successfully"})
}

type transferLoanOfficerRequest struct {
	LoanOfficerID uint32 `binding:"required" json:"loan_officer_id"`
	AdminID       uint32 `binding:"required" json:"admin_id"`
}

func (s *Server) transferLoanOfficer(ctx *gin.Context) {
	var req transferLoanOfficerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if err := s.repo.Loans.TransferLoan(ctx, req.LoanOfficerID, id, req.AdminID); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan officer successfully transferred"})
}

func (s *Server) getLoan(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	loan, err := s.repo.Loans.GetLoanByID(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp, err := s.structureLoan(&loan, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) listLoansByCategory(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	params := repository.Category{}

	b := ctx.Query("branch")
	if b != "" {
		branchID, err := pkg.StringToUint32(b)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.BranchID = pkg.Uint32Ptr(branchID)
	}

	c := ctx.Query("client")
	if c != "" {
		clientID, err := pkg.StringToUint32(c)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.ClientID = pkg.Uint32Ptr(clientID)
	}

	l := ctx.Query("loan_officer")
	if l != "" {
		loanOfficer, err := pkg.StringToUint32(l)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.LoanOfficer = pkg.Uint32Ptr(loanOfficer)
	}

	status := ctx.Query("status")
	if status != "" {
		params.Status = pkg.StringPtr(strings.ToUpper(ctx.Param("status")))
	}

	loans, err := s.repo.Loans.ListLoans(ctx, &params, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]loanResponse, len(loans))

	for idx, l := range loans {
		rsp[idx], err = s.structureLoan(&l, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}

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

	if client.Dob == nil {
		client.Dob = pkg.TimePtr(time.Time{})
	}

	if client.IdNumber == nil {
		client.IdNumber = pkg.StringPtr("")
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

	createdByUser, err := s.repo.Users.GetUserByID(ctx, loan.CreatedBy)
	if err != nil {
		return loanResponse{}, err
	}

	createdBy := loanUserResponse{
		ID:       createdByUser.ID,
		Fullname: createdByUser.FullName,
		Role:     createdByUser.Role,
	}

	if loan.DueDate == nil {
		loan.DueDate = &time.Time{}
	}

	if loan.DisbursedOn == nil {
		loan.DisbursedOn = &time.Time{}
	}

	if loan.LoanPurpose == nil {
		loan.LoanPurpose = pkg.StringPtr("")
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
		FeePaid:            loan.FeePaid,
		PaidAmount:         loan.PaidAmount,
		UpdatedBy:          updatedBy,
		CreatedBy:          loanUserResponse{ID: createdBy.ID, Fullname: createdBy.Fullname, Role: createdBy.Role},
		CreatedAt:          loan.CreatedAt,
	}, nil
}
