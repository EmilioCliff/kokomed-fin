package handlers

import (
	"net/http"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) getDashboardData(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Listing Branches")
	defer span.End()

	data, err := s.repo.Helpers.GetDashboardData(tc)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return 
	}

	ctx.JSON(http.StatusOK, data)
}

func (s *Server) getLoanFormData(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Gettting Helper Loan Form Data")
	defer span.End()

	product := ctx.Query("products")
	var err error
	var products []repository.ProductData
	if product != "" {
		span.SetAttributes(attribute.Bool("products", true))
		products, err = s.repo.Helpers.GetProductData(tc)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var clients []repository.ClientData
	client := ctx.Query("client")
	if client != "" {
		span.SetAttributes(attribute.Bool("clients", true))
		clients, err = s.repo.Helpers.GetClientData(tc)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}
	
	var users []repository.LoanOfficerData
	user := ctx.Query("user")
	if user != "" {
		span.SetAttributes(attribute.Bool("users", true))
		users, err = s.repo.Helpers.GetLoanOfficerData(tc)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var branches []repository.BranchData
	branch := ctx.Query("branch")
	if branch != "" {
		span.SetAttributes(attribute.Bool("branches", true))
		branches, err = s.repo.Helpers.GetBranchData(tc)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
	
			return
		}
	}

	var loans []repository.LoanData
	loan := ctx.Query("loan")
	if loan != "" {
		span.SetAttributes(attribute.Bool("loans", true))
		loans, err = s.repo.Helpers.GetLoanData(tc)
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
	tc, span := s.tracer.Start(ctx.Request.Context(), "Gettting Helper Loan Events")
	defer span.End()

	events, err := s.repo.Loans.GetLoanEvents(tc)
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
	tc, span := s.tracer.Start(ctx.Request.Context(), "Gettting Helper Client Non-Posted")
	defer span.End()

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

	span.SetAttributes(
		attribute.String("page_no", pageNoStr),
		attribute.String("page_size", pageSizeStr),
	)

	rslt, pgData, err := s.repo.Helpers.GetClientNonPayments(tc, req.ID, req.PhoneNumber, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data": rslt,
	})
}