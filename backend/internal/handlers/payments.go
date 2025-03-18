package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) paymentCallback(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Payment Callback")
	defer span.End()

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
		PhoneNumber:   "***",							// safaricom hidden number
		PayingName:    req["FirstName"].(string),
		Amount:        amountFlt,
		AssignedBy:    "APP",
	}

	span.SetAttributes(
		attribute.String("transaction_id", callbackData.TransactionID),
		attribute.String("account_number", callbackData.AccountNumber),
		attribute.String("phone_number", callbackData.PhoneNumber),
		attribute.String("paying_name", callbackData.PayingName),
		attribute.Float64("amount", callbackData.Amount),
		attribute.Bool("assigned", false),
	)

	if app, ok := req["App"].(string); ok && app != "" {
		callbackData.TransactionSource = "INTERNAL"
		
		if email, ok := req["Email"].(string); ok && email != "" {
			callbackData.AssignedBy = email
		}

		span.SetAttributes(
			attribute.String("transaction_source", callbackData.TransactionSource),
			attribute.String("assigned_by", callbackData.AssignedBy),
		)

	} else {
		callbackData.TransactionSource = "MPESA"

		span.SetAttributes(attribute.String("transaction_source", callbackData.TransactionSource))
	}

	if paidDate, ok := req["DatePaid"].(string); ok && paidDate != "" {
		paidDateT, err := time.Parse("2006-01-02", paidDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid disburse_on date format")))

			return
		}

		callbackData.PaidDate = pkg.TimePtr(paidDateT)
		span.SetAttributes(attribute.String("paid_date", paidDate))
	}

	clientID, err := s.repo.Clients.GetClientIDByPhoneNumber(ctx, callbackData.AccountNumber)
	if err != nil {
		if pkg.ErrorCode(err) != pkg.NOT_FOUND_ERROR {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
			
			return
		}
	}
	
	// if account number is wrong(non-existing) continue
	if clientID != 0 {
		span.SetAttributes(
			attribute.Int64("client_id", int64(clientID)),
			attribute.Bool("assigned", true),
		)
		callbackData.AssignedTo = pkg.Uint32Ptr(clientID)
		s.cache.Del(ctx, fmt.Sprintf("client:%v", clientID))
	}

	loanId, err := s.payments.ProcessCallback(tc, &callbackData); 
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
	tc, span := s.tracer.Start(ctx.Request.Context(), "Assigning Payment to Client")
	defer span.End()

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

	span.SetAttributes(
		attribute.Int64("non_posted_id", int64(id)),
		attribute.Int64("client_id", int64(req.ClientID)),
		attribute.String("assigned_by", payloadData.Email),
	)

	loanId, err := s.payments.TriggerManualPayment(tc, services.ManualPaymentData{
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

	accessToken, err := pkg.GenerateAccessToken(s.config.MPESA_CONSUMER_KEY, s.config.MPESA_CONSUMER_SECRET)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (s *Server) validationCallback(ctx *gin.Context) {
	// No validation for my darajaApp
}