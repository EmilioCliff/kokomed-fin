package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) paymentCallback(ctx *gin.Context) {
	var rq any
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"ResultCode": 400,
			"ResultDesc": "Rejected",
		})

		return
	}

	req, _ := rq.(map[string]interface{})

	amountFlt, err := strconv.ParseFloat(req["TransAmount"].(string), 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"ResultCode": 400,
			"ResultDesc": "Rejected",
		})

		return
	}

	callbackData := services.MpesaCallbackData{
		TransactionID: req["TransID"].(string),
		AccountNumber: req["BillRefNumber"].(string),
		PhoneNumber:   "***", // safaricom hidden number
		PayingName:    req["FirstName"].(string),
		Amount:        amountFlt,
		AssignedBy:    "APP",
		AssignedTo:    nil,
	}

	if app, ok := req["App"].(string); ok && app != "" {
		callbackData.TransactionSource = "INTERNAL"

		if email, ok := req["Email"].(string); ok && email != "" {
			callbackData.AssignedBy = email
		}
	} else {
		callbackData.TransactionSource = "MPESA"
	}

	if paidDate, ok := req["DatePaid"].(string); ok && paidDate != "" {
		paidDateT, err := time.Parse("2006-01-02", paidDate)
		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid disburse_on date format")),
			)

			return
		}

		callbackData.PaidDate = pkg.TimePtr(paidDateT)
	}

	clientID, err := s.repo.Clients.GetClientIDByPhoneNumber(ctx, callbackData.AccountNumber)
	if err != nil && pkg.ErrorCode(err) != pkg.NOT_FOUND_ERROR {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	if clientID != 0 {
		callbackData.AssignedTo = pkg.Uint32Ptr(clientID)
		s.cache.Del(ctx, fmt.Sprintf("client:%v", clientID))
	}

	loanId, err := s.payments.ProcessCallback(ctx, &callbackData)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	s.cache.Del(ctx, fmt.Sprintf("loan:%d", loanId))
	s.cache.DelAll(ctx, "non-posted/all:limit=*")

	s.cache.DelAll(ctx, "loan:limit=*")
	s.cache.DelAll(ctx, "client:limit=*")

	ctx.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "Accepted",
	})
}

type paymentByAdminRequest struct {
	ClientID uint32 `binding:"required" json:"clientId"`
}

func (s *Server) paymentByAdmin(ctx *gin.Context) {
	var req paymentByAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

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

	if strings.ToLower(payloadData.Role) != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})

		return
	}

	loanId, err := s.payments.TriggerManualPayment(ctx, services.ManualPaymentData{
		NonPostedID: id,
		ClientID:    req.ClientID,
		AdminUserID: payloadData.UserID,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	s.cache.Del(ctx, fmt.Sprintf("loan:%d", loanId))
	s.cache.DelAll(ctx, "loan:limit=*")

	s.cache.Del(ctx, fmt.Sprintf("non-posted:%d", id))
	s.cache.DelAll(ctx, "non-posted/all:limit=*")

	s.cache.Del(ctx, fmt.Sprintf("client:%v", req.ClientID))
	s.cache.DelAll(ctx, "client:limit=*")

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) updateLoan(ctx *gin.Context) {}

func (s *Server) deleteLoan(ctx *gin.Context) {}

func (s *Server) getMPESAAccesToken(ctx *gin.Context) {
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

	if payloadData.Email != "emiliocliff@gmail.com" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})

		return
	}

	accessToken, err := pkg.GenerateAccessToken(
		s.config.MPESA_CONSUMER_KEY,
		s.config.MPESA_CONSUMER_SECRET,
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (s *Server) validationCallback(ctx *gin.Context) {
	// No validation for my darajaApp
}
