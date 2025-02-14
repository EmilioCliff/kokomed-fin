package handlers

import (
	"context"
	"fmt"
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "%v", err.Error())))
		return
	}

	var err error
	filters := services.ReportFilters{}
	if req.StartDate != "" {
		if filters.StartDate, err = time.Parse("2006-01-02", req.StartDate); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid start date format")))
			return
		}
	}
	if req.EndDate != "" {
		if filters.EndDate, err = time.Parse("2006-01-02", req.EndDate); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid end date format")))
			return
		}
	}

	format := ctx.Query("format")
	fileExt := "pdf" 

	baseName := req.ReportName + "_report"
	if req.UserId != 0 {
		filters.UserId = pkg.Uint32Ptr(req.UserId)
		baseName = fmt.Sprintf("user_%d", req.UserId)
	} else if req.ClientId != 0 {
		filters.ClientId = pkg.Uint32Ptr(req.ClientId)
		baseName = fmt.Sprintf("client_%d", req.ClientId)
	} else if req.LoanId != 0 {
		filters.LoanId = pkg.Uint32Ptr(req.LoanId)
		baseName = fmt.Sprintf("loan_LN%03d", req.LoanId)
	} else {
		if format == "excel" {
			fileExt = "xlsx"
		} else if format != "pdf" {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "valid format is required (pdf or excel)")))
			return
		}
	}

	fileName := fmt.Sprintf("%s_%s.%s", baseName, time.Now().Format("2006-01-02"), fileExt)

	reportGenerators := map[string]func(context.Context, string, services.ReportFilters) ([]byte, error){
		"payment": s.report.GeneratePaymentsReport,
		"branch":  s.report.GenerateBranchesReport,
		"user":    s.report.GenerateUsersReport,
		"client":  s.report.GenerateClientsReport,
		"product": s.report.GenerateProductsReport,
		"loan":    s.report.GenerateLoansReport,
	}

	generateReport, exists := reportGenerators[req.ReportName]
	if !exists {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid report name")))
		return
	}

	responseBytes, reportErr := generateReport(ctx, format, filters)
	if reportErr != nil || len(responseBytes) == 0 {
		ctx.JSON(http.StatusInternalServerError, errorResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "report generation failed")))
		return
	}

	contentType := map[string]string{
		"pdf":  "application/pdf",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}[fileExt]

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	ctx.Data(http.StatusOK, contentType, responseBytes)
}