package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type downloadReportsRequest struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	ReportName string `json:"reportName" binding:"required"`
	UserId uint32 `json:"userId"`
	ClientId uint32 `json:"clientId"`
	LoanId uint32 `json:"loanId"`
}


func (s *Server) generateReport(ctx *gin.Context) {
	var req downloadReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	fromDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid date format")))

		return
	}

	toDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid date format")))

		return
	}

	format := ctx.Query("format")
	if format == "" || (format != "pdf" && format != "excel") {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "valid format is required")))
		
		return
	}

	filters := services.ReportFilters{
		StartDate: fromDate,
		EndDate:   toDate,
		UserId: nil,
		ClientId: nil,
	}

	if req.UserId != 0 {
		filters.UserId = pkg.Uint32Ptr(req.UserId)
	}

	if req.ClientId != 0 {
		filters.ClientId = pkg.Uint32Ptr(req.ClientId)
	}

	if req.LoanId != 0 {
		filters.LoanId = pkg.Uint32Ptr(req.LoanId)
	}

	var reportErr error
	switch req.ReportName{
		case "payment":
			reportErr = s.report.GeneratePaymentsReport(ctx, format, filters)
			
		case "branch":
			reportErr = s.report.GenerateBranchesReport(ctx, format, filters)

		case "user":
			reportErr = s.report.GenerateUsersReport(ctx, format, filters)
		
		case "client":
			reportErr = s.report.GenerateClientsReport(ctx, format, filters)

		case "product":
			reportErr = s.report.GenerateProductsReport(ctx, format, filters)

		case "loan":
			reportErr = s.report.GenerateLoansReport(ctx, format, filters)

		default:
			ctx.JSON(http.StatusBadRequest, pkg.Errorf(pkg.INVALID_ERROR, "invalid report name"))

			return
	}
	if reportErr != nil {
		ctx.JSON(pkg.ErrorToStatusCode(reportErr), errorResponse(reportErr))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "report generated successfully"})

	// data, err := s.report.GeneratePaymentsReport(format, services.ReportFilters{
	// 	StartDate: fromDate,
	// 	EndDate:   toDate,
	// })
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// ctx.Data(http.StatusOK, "application/octet-stream", data)
}