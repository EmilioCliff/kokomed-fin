package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
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
	RemainingAmount float64 `json:"remainingAmount"`
	UpdatedBy          userResponse	`json:"updatedBy"`
	CreatedBy          userResponse	`json:"createdBy"`
	CreatedAt          time.Time	`json:"createdAt"`
}

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
	tc, span := s.tracer.Start(ctx.Request.Context(), "Creating Loan")
	defer span.End()

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
	span.SetAttributes(attribute.String("createdBy", payloadData.Email))

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
		Status: 		   "INACTIVE",
	}

	if req.ProcessingFeePaid {
		params.FeePaid = true
		span.SetAttributes(attribute.Bool("processingFeePaid", true))
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
		params.Status = "ACTIVE"
		params.DueDate = pkg.TimePtr(disburseDate.AddDate(0, 0, int(req.Installments)*int(req.InstallmentsPeriod)))
		
		span.SetAttributes(
			attribute.String("disbursedDate", disburseDate.String()),
			attribute.String("dueDate", disburseDate.AddDate(0, 0, int(req.Installments)*int(req.InstallmentsPeriod)).String()),
		)
	}

	if req.LoanPurpose != "" {
		params.LoanPurpose = pkg.StringPtr(req.LoanPurpose)
	}

	loan, err := s.repo.Loans.CreateLoan(tc, &params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// rsp, err := s.structureLoan(&loan, ctx)
	// if err != nil {
	// 	ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

	// 	return
	// }

	s.cache.DelAll(ctx, "loan:limit=*")
	
	s.cache.DelAll(ctx, "client:limit=*")
	s.cache.Del(ctx, fmt.Sprintf("client:%v", req.ClientID))

	ctx.JSON(http.StatusOK, loan)
}

// binding:"oneof=ACTIVE DEFAULTED"
type disburseLoanRequest struct {
	Status string ` json:"status"`
	DisburseDate string `json:"disburseDate"`
	FeePaid bool `json:"feePaid"`
}

func (s *Server) disburseLoan(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Disbursing Loan")
	defer span.End()

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil || id == 0 {
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

	span.SetAttributes(
		attribute.String("disbursedBy", payloadData.Email),
		attribute.String("loanID", fmt.Sprint(id)),
		attribute.String("status", req.Status),
		attribute.Bool("feePaid", req.FeePaid),
	)

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

	clientId, err := s.repo.Loans.DisburseLoan(tc, &params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	s.cache.Del(ctx, fmt.Sprintf("loan:%d", id))
	s.cache.DelAll(ctx, "loan:limit=*")
	
	s.cache.Del(ctx, fmt.Sprintf("client:%v", clientId))
	s.cache.DelAll(ctx, "client:limit=*")

	ctx.JSON(http.StatusOK, gin.H{"success": "Loan disbursed successfully"})
}

func (s *Server) getLoanInstallments(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Getting Loan Installments")
	defer span.End()

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	installments, err := s.repo.Loans.GetLoanInstallments(tc, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": installments})
}

func (s *Server) listLoansByCategory(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Listing Loans")
	defer span.End()
	
	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	span.SetAttributes(
		attribute.String("page_no", pageNoStr),
		attribute.String("page_size", pageSizeStr),
	)

	params := repository.Category{}
	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
	}

	b := ctx.Query("branch")
	if b != "" {
		branchID, err := pkg.StringToUint32(b)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		span.SetAttributes(attribute.String("branch", b))
		params.BranchID = pkg.Uint32Ptr(branchID)
		cacheParams["branch"] = []string{b}
	}

	search := ctx.Query("search")
	if search != "" {
		span.SetAttributes(attribute.String("searched", search))
		params.Search = pkg.StringPtr(strings.ToLower(search))
		cacheParams["search"] = []string{search}
	}

	status := ctx.Query("status")
	if status != "" {
		statuses := strings.Split(status, ",")
		
		for i := range statuses {
			statuses[i] = strings.TrimSpace(statuses[i])
		}
		
		span.SetAttributes(attribute.String("status", status))
		params.Statuses = pkg.StringPtr(strings.Join(statuses, ","))
		cacheParams["status"] = []string{strings.Join(statuses, ",")}
	}

	loans, pgData, err := s.repo.Loans.ListLoans(tc, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// rsp := make([]loanResponse, len(loans))

	// for idx, l := range loans {
	// 	rsp[idx], err = s.structureLoan(&l, ctx)
	// 	if err != nil {
	// 		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

	// 		return
	// 	}
	// }

	response := gin.H{
		"metadata": pgData,
		"data": loans, 
	}

	cacheKey := constructCacheKey("loan", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) listUnpaidInstallmentsData(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Listing Unpaid Loans Installments")
	defer span.End()

	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	span.SetAttributes(
		attribute.String("page_no", pageNoStr),
		attribute.String("page_size", pageSizeStr),
	)

	params := repository.Category{}
	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
	}

	search := ctx.Query("search")
	if search != "" {
		span.SetAttributes(attribute.String("searched", search))
		params.Search = pkg.StringPtr(strings.ToLower(search))
		cacheParams["search"] = []string{search}
	}

	rsp, pgData, err := s.repo.Loans.ListUnpaidInstallmentsData(tc, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	response := gin.H{
		"metadata": pgData,
		"data": rsp, 
	}

	cacheKey := constructCacheKey("loan/unpaid-installments", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) listExpectedPayments(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Listing Expected Payment")
	defer span.End()

	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	span.SetAttributes(
		attribute.String("page_no", pageNoStr),
		attribute.String("page_size", pageSizeStr),
	)

	params := repository.Category{}
	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
	}

	search := ctx.Query("search")
	if search != "" {
		span.SetAttributes(attribute.String("searched", search))
		params.Search = pkg.StringPtr(strings.ToLower(search))
		cacheParams["search"] = []string{search}
	}

	expectedPayments, pgData, err := s.repo.Loans.GetExpectedPayments(tc, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	response := gin.H{
		"metadata": pgData,
		"data": expectedPayments,
	}

	cacheKey := constructCacheKey("loan/expected-payments", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
}

// func (s *Server) // func (s *Server) structureLoan(loan *repository.Loan, ctx *gin.Context) (loanResponse, error) {
// 	cacheKey := fmt.Sprintf("loan:%v", loan.ID)
// 	var dataCached loanResponse

// 	exists, _ := s.cache.Get(ctx, cacheKey, &dataCached)
// 	if exists {
// 		return dataCached, nil
// 	}

// 	product, err := s.repo.Products.GetProductByID(ctx, loan.ProductID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	productBranch, err := s.repo.Branches.GetBranchByID(ctx, product.BranchID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	client, err := s.repo.Clients.GetClient(ctx, loan.ClientID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	if client.Dob == nil {
// 		client.Dob = pkg.TimePtr(time.Time{})
// 	}

// 	if client.IdNumber == nil {
// 		client.IdNumber = pkg.StringPtr("")
// 	}

// 	clientBranch, err := s.repo.Branches.GetBranchByID(ctx, client.BranchID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	loanOfficer, err := s.repo.Users.GetUserByID(ctx, loan.LoanOfficerID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	approvedBy, err := s.repo.Users.GetUserByID(ctx, loan.ApprovedBy)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	disbursedBy := userResponse{}
	
// 	if loan.DisbursedBy != nil {
// 		disbursedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.DisbursedBy)
// 		if err != nil {
// 			return loanResponse{}, err
// 		}

// 		disbursedBy = userResponse{
// 			ID:       disbursedByUser.ID,
// 			Fullname: disbursedByUser.FullName,
// 			Email: disbursedByUser.Email,
// 			PhoneNumber: disbursedByUser.PhoneNumber,
// 		}
// 	}

// 	updatedBy := userResponse{}

// 	if loan.UpdatedBy != nil {
// 		updatedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.UpdatedBy)
// 		if err != nil {
// 			return loanResponse{}, err
// 		}

// 		updatedBy = userResponse{
// 			ID:       updatedByUser.ID,
// 			Fullname: updatedByUser.FullName,
// 			Email: updatedByUser.Email,
// 			PhoneNumber: updatedBy.PhoneNumber,
// 		}
// 	}

// 	createdByUser, err := s.repo.Users.GetUserByID(ctx, loan.CreatedBy)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	createdBy := userResponse{
// 		ID:       createdByUser.ID,
// 		Fullname: createdByUser.FullName,
// 		Email: createdByUser.Email,
// 		PhoneNumber: createdByUser.PhoneNumber,
// 	}

// 	if loan.DueDate == nil {
// 		loan.DueDate = &time.Time{}
// 	}

// 	if loan.DisbursedOn == nil {
// 		loan.DisbursedOn = &time.Time{}
// 	}

// 	if loan.LoanPurpose == nil {
// 		loan.LoanPurpose = pkg.StringPtr("")
// 	}

// 	rsp := loanResponse{
// 		ID: loan.ID,
// 		Product: productResponse{
// 			ID:             product.ID,
// 			BranchName:     productBranch.Name,
// 			LoanAmount:     product.LoanAmount,
// 			RepayAmount:    product.RepayAmount,
// 			InterestAmount: product.InterestAmount,
// 		},
// 		Client: clientResponse{
// 			ID:          client.ID,
// 			FullName:        client.FullName,
// 			PhoneNumber: client.PhoneNumber,
// 			Active:      client.Active,
// 			BranchName:  clientBranch.Name,
// 		},
// 		LoanOfficer:        userResponse{
// 			ID:       loanOfficer.ID,
// 			Fullname: loanOfficer.FullName,
// 			Email: loanOfficer.Email,
// 			PhoneNumber: loanOfficer.PhoneNumber,
// 		},
// 		LoanPurpose:        *loan.LoanPurpose,
// 		DueDate:            *loan.DueDate,
// 		ApprovedBy:         userResponse{
// 			ID:       approvedBy.ID,
// 			Fullname: approvedBy.FullName,
// 			Email: approvedBy.Email,
// 			PhoneNumber: approvedBy.PhoneNumber,
// 		},
// 		DisbursedOn:        *loan.DisbursedOn,
// 		DisbursedBy:        disbursedBy,
// 		TotalInstallments:  loan.TotalInstallments,
// 		InstallmentsPeriod: loan.InstallmentsPeriod,
// 		Status:             loan.Status,
// 		ProcessingFee:      loan.ProcessingFee,
// 		FeePaid:            loan.FeePaid,
// 		PaidAmount:         loan.PaidAmount,
// 		RemainingAmount: 	product.RepayAmount - loan.PaidAmount,
// 		UpdatedBy:          updatedBy,
// 		CreatedBy:          createdBy,
// 		CreatedAt:          loan.CreatedAt,
// 	}

// 	if err := s.cache.Set(ctx, cacheKey, rsp, 3*time.Minute); err != nil {
// 		return loanResponse{}, err
// 	}

// 	return rsp, nil
// }(loan *repository.Loan, ctx *gin.Context) (loanResponse, error) {
// 	cacheKey := fmt.Sprintf("loan:%v", loan.ID)
// 	var dataCached loanResponse

// 	exists, _ := s.cache.Get(ctx, cacheKey, &dataCached)
// 	if exists {
// 		return dataCached, nil
// 	}

// 	product, err := s.repo.Products.GetProductByID(ctx, loan.ProductID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	productBranch, err := s.repo.Branches.GetBranchByID(ctx, product.BranchID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	client, err := s.repo.Clients.GetClient(ctx, loan.ClientID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	if client.Dob == nil {
// 		client.Dob = pkg.TimePtr(time.Time{})
// 	}

// 	if client.IdNumber == nil {
// 		client.IdNumber = pkg.StringPtr("")
// 	}

// 	clientBranch, err := s.repo.Branches.GetBranchByID(ctx, client.BranchID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	loanOfficer, err := s.repo.Users.GetUserByID(ctx, loan.LoanOfficerID)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	approvedBy, err := s.repo.Users.GetUserByID(ctx, loan.ApprovedBy)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	disbursedBy := userResponse{}
	
// 	if loan.DisbursedBy != nil {
// 		disbursedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.DisbursedBy)
// 		if err != nil {
// 			return loanResponse{}, err
// 		}

// 		disbursedBy = userResponse{
// 			ID:       disbursedByUser.ID,
// 			Fullname: disbursedByUser.FullName,
// 			Email: disbursedByUser.Email,
// 			PhoneNumber: disbursedByUser.PhoneNumber,
// 		}
// 	}

// 	updatedBy := userResponse{}

// 	if loan.UpdatedBy != nil {
// 		updatedByUser, err := s.repo.Users.GetUserByID(ctx, *loan.UpdatedBy)
// 		if err != nil {
// 			return loanResponse{}, err
// 		}

// 		updatedBy = userResponse{
// 			ID:       updatedByUser.ID,
// 			Fullname: updatedByUser.FullName,
// 			Email: updatedByUser.Email,
// 			PhoneNumber: updatedBy.PhoneNumber,
// 		}
// 	}

// 	createdByUser, err := s.repo.Users.GetUserByID(ctx, loan.CreatedBy)
// 	if err != nil {
// 		return loanResponse{}, err
// 	}

// 	createdBy := userResponse{
// 		ID:       createdByUser.ID,
// 		Fullname: createdByUser.FullName,
// 		Email: createdByUser.Email,
// 		PhoneNumber: createdByUser.PhoneNumber,
// 	}

// 	if loan.DueDate == nil {
// 		loan.DueDate = &time.Time{}
// 	}

// 	if loan.DisbursedOn == nil {
// 		loan.DisbursedOn = &time.Time{}
// 	}

// 	if loan.LoanPurpose == nil {
// 		loan.LoanPurpose = pkg.StringPtr("")
// 	}

// 	rsp := loanResponse{
// 		ID: loan.ID,
// 		Product: productResponse{
// 			ID:             product.ID,
// 			BranchName:     productBranch.Name,
// 			LoanAmount:     product.LoanAmount,
// 			RepayAmount:    product.RepayAmount,
// 			InterestAmount: product.InterestAmount,
// 		},
// 		Client: clientResponse{
// 			ID:          client.ID,
// 			FullName:        client.FullName,
// 			PhoneNumber: client.PhoneNumber,
// 			Active:      client.Active,
// 			BranchName:  clientBranch.Name,
// 		},
// 		LoanOfficer:        userResponse{
// 			ID:       loanOfficer.ID,
// 			Fullname: loanOfficer.FullName,
// 			Email: loanOfficer.Email,
// 			PhoneNumber: loanOfficer.PhoneNumber,
// 		},
// 		LoanPurpose:        *loan.LoanPurpose,
// 		DueDate:            *loan.DueDate,
// 		ApprovedBy:         userResponse{
// 			ID:       approvedBy.ID,
// 			Fullname: approvedBy.FullName,
// 			Email: approvedBy.Email,
// 			PhoneNumber: approvedBy.PhoneNumber,
// 		},
// 		DisbursedOn:        *loan.DisbursedOn,
// 		DisbursedBy:        disbursedBy,
// 		TotalInstallments:  loan.TotalInstallments,
// 		InstallmentsPeriod: loan.InstallmentsPeriod,
// 		Status:             loan.Status,
// 		ProcessingFee:      loan.ProcessingFee,
// 		FeePaid:            loan.FeePaid,
// 		PaidAmount:         loan.PaidAmount,
// 		RemainingAmount: 	product.RepayAmount - loan.PaidAmount,
// 		UpdatedBy:          updatedBy,
// 		CreatedBy:          createdBy,
// 		CreatedAt:          loan.CreatedAt,
// 	}

// 	if err := s.cache.Set(ctx, cacheKey, rsp, 3*time.Minute); err != nil {
// 		return loanResponse{}, err
// 	}

// 	return rsp, nil
// }
