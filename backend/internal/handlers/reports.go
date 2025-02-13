package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type downloadReportsRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	ReportName string `json:"reportName"`
	UserId uint32 `json:"userId"`
	ClientId uint32 `json:"clientId"`
	LoanId uint32 `json:"loanId"`
}


func (s *Server) generateReport(ctx *gin.Context) {
	var req downloadReportsRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	filters := services.ReportFilters{
		UserId: nil,
		ClientId: nil,
		LoanId: nil,
	}

	if req.StartDate != "" {
		filters.StartDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid date format")))
	
			return
		}
	}

	if req.EndDate != "" {
		filters.EndDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid date format")))
	
			return
		}
	}

	var contentType string
	var fileExt string
	format := ctx.Query("format")
	if format == "excel" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		fileExt = "xlsx"
	} else if format == "pdf" {
		contentType = "application/pdf"
		fileExt = "pdf"
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "valid format is required (pdf or excel)")))

		return
	}

	fileName := "report" 

	if req.UserId != 0 {
		filters.UserId = pkg.Uint32Ptr(req.UserId)
		fileName = fmt.Sprintf("user_%d", req.UserId)
	} else if req.ClientId != 0 {
		filters.ClientId = pkg.Uint32Ptr(req.ClientId)
		fileName = fmt.Sprintf("client_%d", req.ClientId)
	} else if req.LoanId != 0 {
		filters.LoanId = pkg.Uint32Ptr(req.LoanId)
		fileName = fmt.Sprintf("loan_LN%03d", req.LoanId)
	}

	var reportErr error
	var responseBytes []byte
	switch req.ReportName{
		case "payment":
			responseBytes, reportErr = s.report.GeneratePaymentsReport(ctx, format, filters)
			fileName = "payment"
			
		case "branch":
			responseBytes, reportErr = s.report.GenerateBranchesReport(ctx, format, filters)
			fileName = "branch"

		case "user":
			reportErr = s.report.GenerateUsersReport(ctx, format, filters)
		
		case "client":
			reportErr = s.report.GenerateClientsReport(ctx, format, filters)

		case "product":
			responseBytes, reportErr = s.report.GenerateProductsReport(ctx, format, filters)
			fileName = "product"

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

	log.Println(contentType)
	log.Println(fileExt)
	log.Println(fileName)
	log.Println(responseBytes)

	ctx.JSON(http.StatusOK, gin.H{"success": "report generated successfully"})

	// if len(responseBytes) == 0 {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "report generation failed")))

	// 	return
	// }

	// ctx.Header("Content-Type", contentType)
	// ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s_%s.%s"`, fileName, time.Now().Format("2006-01-02"), fileExt))

	// ctx.Data(http.StatusOK, contentType, responseBytes)
}