package handlers

import (
	"net/http"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) getDashboardData(ctx *gin.Context) {
	// get inactiveLoans

	// get recentPayments

	// get widgetData
}

func (s *Server) getLoanFormData(ctx *gin.Context) {
	// get allLoans
	products, err := s.repo.Helpers.GetProductData()
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	clients, err := s.repo.Helpers.GetClientData()
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}
	
	loanOfficers, err := s.repo.Helpers.GetLoanOfficerData()
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":     products,
		"clients":      clients,
		"loanOfficers": loanOfficers,
})
}