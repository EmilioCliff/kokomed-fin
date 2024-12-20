package handlers

import (
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
		log.Println("failed to bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	req, _ := rq.(map[string]interface{})

	amount, err := strconv.ParseFloat(req["TransAmount"].(string), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	callbackData := services.MpesaCallbackData{
		TransactionID: req["TransID"].(string),
		AccountNumber: req["BillRefNumber"].(string),
		PhoneNumber:   req["MSISDN"].(string),
		PayingName:    req["FirstName"].(string),
		Amount:        amount,
	}

	clientID, err := s.repo.Clients.GetClientIDByPhoneNumber(ctx, callbackData.AccountNumber)
	if err != nil {
		// if account number is wrong(non-existing) continue
		if pkg.ErrorCode(err) != pkg.NOT_FOUND_ERROR {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}
	}

	if clientID != 0 {
		callbackData.AssignedTo = pkg.Uint32Ptr(clientID)
	}

	if err := s.payments.ProcessCallback(ctx, &callbackData); err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "Accepted",
	})
}

type paymentByAdminRequest struct {
	ClientID uint32 `binding:"required" json:"client_id"`
	AdminID  uint32 `binding:"required" json:"admin_id"`
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

	err = s.payments.TriggerManualPayment(ctx, services.ManualPaymentData{
		NonPostedID: id,
		ClientID:    req.ClientID,
		AdminUserID: req.AdminID,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
