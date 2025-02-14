package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		log.Println(err)
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

	if app, ok := req["App"].(string); ok && app != "" {
		callbackData.TransactionSource = "INTERNAL"
		// payload, ok := ctx.Get(authorizationPayloadKey)
		// if !ok {
		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})

		// 	return
		// }

		// payloadData, ok := payload.(*pkg.Payload)
		// if !ok {
		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect token"})

		// 	return
		// }

		callbackData.AssignedBy = "emiliocliff@gmail.com"
	} else {
		callbackData.TransactionSource = "MPESA"
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
		callbackData.AssignedTo = pkg.Uint32Ptr(clientID)
	}

	loanId, err := s.payments.ProcessCallback(ctx, &callbackData); 
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	s.cache.Del(ctx, fmt.Sprintf("loan:%d", loanId))
	s.cache.DelAll(ctx, "non-posted/all:limit=*")

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
	s.cache.Del(ctx, fmt.Sprintf("non-posted:%d", id))
	s.cache.DelAll(ctx, "non-posted/all:limit=*")

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