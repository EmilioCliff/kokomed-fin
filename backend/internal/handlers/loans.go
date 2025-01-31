package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type loanResponse struct {
	ID                 uint32	`json:"id"`
	Product            productResponse	`json:"product"`
	Client             clientResponse	`json:"client"`
	LoanOfficer        userResponse	`json:"loanOfficer"`
	LoanPurpose        string	`json:"loanPurpose"`
	DueDate            time.Time	`json:"dueDate"`
	ApprovedBy         userResponse	`json:"approvedBy"`
	DisbursedOn        time.Time	`json:"disbursedOn"`
	DisbursedBy        userResponse	`json:"disbursedBy"`
	TotalInstallments  uint32	`json:"totalInstallments"`
	InstallmentsPeriod uint32	`json:"installmentsPeriod"`
	Status             string	`json:"status"`
	ProcessingFee      float64	`json:"processingFee"`
	FeePaid            bool	`json:"feePaid"`
	PaidAmount         float64	`json:"paidAmount"`
	UpdatedBy          userResponse	`json:"updatedBy"`
	CreatedBy          userResponse	`json:"createdBy"`
	CreatedAt          time.Time	`json:"createdAt"`
}

// type loanProductResponse struct {
// 	ID             uint32  `json:"id"`
// 	BranchName     string  `json:"branch_name"`
// 	LoanAmount     float64 `json:"loan_amount"`
// 	RepayAmount    float64 `json:"repay_amount"`
// 	InterestAmount float64 `json:"interest_amount"`
// }

// type loanClientResponse struct {
// 	ID          uint32    `json:"id"`
// 	Name        string    `json:"name"`
// 	PhoneNumber string    `json:"phone_number"`
// 	IdNumber    string    `json:"id_number"`
// 	Dob         time.Time `json:"dob"`
// 	Gender      string    `json:"gender"`
// 	Active      bool      `json:"active"`
// 	BranchName  string    `json:"branch_name"`
// 	Overpayment float64   `json:"overpayment"`
// 	AssignedStaff userResponse `json:"assigned_staff"`
// 	DueAmount float64 `json:"due_amount"`
// 	CreatedBy userResponse `json:"created_by"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type loanUserResponse struct {
// 	ID       uint32 `json:"id"`
// 	Fullname string `json:"fullname"`
// 	PhoneNumber string `json:"phone_number"`
// 	Email string `json:"email"`
// 	Role     string `json:"role"`
// 	BranchName string `json:"branch_name"`
// 	CreatedAt time.Time `json:"created_at"`
// }

type createLoanRequest struct {
	ProductID          uint32  `binding:"required" json:"productId"`
	ClientID           uint32  `binding:"required" json:"clientId"`
	LoanOfficerID      uint32  `binding:"required" json:"loanOfficerId"`
	LoanPurpose        string  `                   json:"loanPurpose"`
	Disburse		   bool    `				   json:"disburse"`
	DisburseOn         string  `                   json:"disburseOn"`
	Installments       uint32  `binding:"required" json:"installments"`
	InstallmentsPeriod uint32  `binding:"required" json:"installmentsPeriod"`
	ProcessingFee      float64 `binding:"required" json:"processingFee"`
	ProcessingFeePaid  bool    `				   json:"processingFeePaid"`
}

func (s *Server) createLoan(ctx *gin.Context) {
	var req createLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	params := repository.Loan{
		ProductID:          req.ProductID,
		ClientID:           req.ClientID,
		LoanOfficerID:      req.LoanOfficerID,
		ApprovedBy:         payloadData.UserID,
		TotalInstallments:  req.Installments,
		InstallmentsPeriod: req.InstallmentsPeriod,
		ProcessingFee:      req.ProcessingFee,
		FeePaid:            false,
		CreatedBy:          payloadData.UserID,
	}

	if req.ProcessingFeePaid {
		params.FeePaid = true
	}

	if req.Disburse {
		var err error
		disburseDate := time.Now()
		if req.DisburseOn != "" {
			disburseDate, err = time.Parse("2006-01-02", req.DisburseOn)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid disburse_on date format")))
	
				return
			}
		} 
		params.DisbursedOn = pkg.TimePtr(disburseDate)
		params.DisbursedBy = pkg.Uint32Ptr(payloadData.UserID)

		params.DueDate = pkg.TimePtr(disburseDate.AddDate(0, 0, int(req.Installments)*int(req.InstallmentsPeriod)))
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
// binding:"oneof=ACTIVE DEFAULTED"
type disburseLoanRequest struct {
	Status string ` json:"status"`
	DisburseDate string `json:"disburseDate"`
	FeePaid bool `json:"feePaid"`
}

func (s *Server) disburseLoan(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		return
	}

	payloadData, ok := payload.(*pkg.Payload)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		return
	}

	var req disburseLoanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	params := repository.DisburseLoan{
		ID:          id,
		DisbursedBy: payloadData.UserID,
	}

	// if status has changed to ACTIVE get the disburse date
	var disbursedDate time.Time
	if req.Status == "ACTIVE" {
		disbursedDate, err = time.Parse("2006-01-02", req.DisburseDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.DisbursedOn = pkg.TimePtr(disbursedDate)
	}

	if req.Status != "" {
		params.Status = pkg.StringPtr(req.Status)
	}

	if req.FeePaid {
		params.FeePaid = &req.FeePaid
	}

	log.Println(req)

	err = s.repo.Loans.DisburseLoan(ctx, &params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Loan disbursed successfully"})
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

	pageSize, err := pkg.StringToUint32(ctx.DefaultQuery("limit", "10"))
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

	search := ctx.Query("search")
	if search != "" {
		params.Search = pkg.StringPtr(strings.ToLower(search))
	}

	status := ctx.Query("status")
	if status != "" {
		statuses := strings.Split(status, ",")
		
		for i := range statuses {
			statuses[i] = strings.TrimSpace(statuses[i])
		}
		
		params.Statuses = pkg.StringPtr(strings.Join(statuses, ","))
	}

	loans, pgData, err := s.repo.Loans.ListLoans(ctx, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
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

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data": rsp, 
	})
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

	disbursedBy := userResponse{}
	
	if loan.DisbursedBy != nil {
		disbursedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.DisbursedBy)
		if err != nil {
			return loanResponse{}, err
		}

		disbursedBy = userResponse{
			ID:       disbursedByUser.ID,
			Fullname: disbursedByUser.FullName,
			Email: disbursedByUser.Email,
			PhoneNumber: disbursedByUser.PhoneNumber,
		}
	}

	updatedBy := userResponse{}

	if loan.UpdatedBy != nil {
		updatedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.UpdatedBy)
		if err != nil {
			return loanResponse{}, err
		}

		updatedBy = userResponse{
			ID:       updatedByUser.ID,
			Fullname: updatedByUser.FullName,
			Email: updatedByUser.Email,
			PhoneNumber: updatedBy.PhoneNumber,
		}
	}

	createdByUser, err := s.repo.Users.GetUserByID(ctx, loan.CreatedBy)
	if err != nil {
		return loanResponse{}, err
	}

	createdBy := userResponse{
		ID:       createdByUser.ID,
		Fullname: createdByUser.FullName,
		Email: createdByUser.Email,
		PhoneNumber: createdByUser.PhoneNumber,
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
		Product: productResponse{
			ID:             product.ID,
			BranchName:     productBranch.Name,
			LoanAmount:     product.LoanAmount,
			RepayAmount:    product.RepayAmount,
			InterestAmount: product.InterestAmount,
		},
		Client: clientResponse{
			ID:          client.ID,
			FullName:        client.FullName,
			PhoneNumber: client.PhoneNumber,
			Active:      client.Active,
			BranchName:  clientBranch.Name,
		},
		LoanOfficer:        userResponse{
			ID:       loanOfficer.ID,
			Fullname: loanOfficer.FullName,
			Email: loanOfficer.Email,
			PhoneNumber: loanOfficer.PhoneNumber,
		},
		LoanPurpose:        *loan.LoanPurpose,
		DueDate:            *loan.DueDate,
		ApprovedBy:         userResponse{
			ID:       approvedBy.ID,
			Fullname: approvedBy.FullName,
			Email: approvedBy.Email,
			PhoneNumber: approvedBy.PhoneNumber,
		},
		DisbursedOn:        *loan.DisbursedOn,
		DisbursedBy:        disbursedBy,
		TotalInstallments:  loan.TotalInstallments,
		InstallmentsPeriod: loan.InstallmentsPeriod,
		Status:             loan.Status,
		ProcessingFee:      loan.ProcessingFee,
		FeePaid:            loan.FeePaid,
		PaidAmount:         loan.PaidAmount,
		UpdatedBy:          updatedBy,
		CreatedBy:          createdBy,
		CreatedAt:          loan.CreatedAt,
	}, nil
}
