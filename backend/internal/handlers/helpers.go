package handlers

import (
	"net/http"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) getDashboardData(ctx *gin.Context) {
	// get inactiveLoans
	data, err := s.repo.Helpers.GetDashboardData(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return 
	}

	ctx.JSON(http.StatusOK, data)
}

func (s *Server) getLoanFormData(ctx *gin.Context) {
	product := ctx.Query("products")
	var err error
	var products []repository.ProductData
	if product != "" {
		products, err = s.repo.Helpers.GetProductData(ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var clients []repository.ClientData
	client := ctx.Query("client")
	if client != "" {
		clients, err = s.repo.Helpers.GetClientData(ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}
	
	var users []repository.LoanOfficerData
	user := ctx.Query("user")
	if user != "" {
		users, err = s.repo.Helpers.GetLoanOfficerData(ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var branches []repository.BranchData
	branch := ctx.Query("branch")
	if branch != "" {
		branches, err = s.repo.Helpers.GetBranchData(ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var loans []repository.LoanData
	loan := ctx.Query("loan")
	if loan != "" {
		loans, err = s.repo.Helpers.GetLoanData(ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product":     products,
		"client":      clients,
		"user": users,
		"branch": branches,
		"loan": loans,
})
}

func (s *Server) getLoanEvents(ctx *gin.Context) {
	events, err := s.repo.Loans.GetLoanEvents(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": events,
	})
}

type getClientNonPostedReq struct {
	ID uint32 `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
}

func (s *Server) getClientNonPosted(ctx *gin.Context) {
	var req getClientNonPostedReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

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

	rslt, pgData, err := s.repo.Helpers.GetClientNonPayments(ctx, req.ID, req.PhoneNumber, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data": rslt,
	})
}