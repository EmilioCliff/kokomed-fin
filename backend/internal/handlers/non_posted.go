package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type nonPostedResponse struct {
	ID                uint32            `json:"id"`
	TransactionSource string            `json:"transactionSource"`
	TransactionNumber string            `json:"transactionNumber"`
	AccountNumber     string            `json:"accountNumber"`
	PhoneNumber       string            `json:"phoneNumber"`
	PayingName        string            `json:"payingName"`
	Amount            float64           `json:"amount"`
	PaidDate          time.Time         `json:"paidDate"`
	AssignedTo        clientShortResponse `json:"assignedTo"`
	Assigned bool	`json:"assigned"`
}

func (s *Server) listAllNonPostedPayments(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSize, err := pkg.StringToUint32(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	params := repository.NonPostedCategory{}

	search := ctx.Query("search")
	if search != "" {
		params.Search = pkg.StringPtr(strings.ToLower(search))
	}

	source := ctx.Query("source")
	if source != "" {
		sources := strings.Split(source, ",")
		
		for i := range sources {
			sources[i] = strings.TrimSpace(sources[i])
		}
		
		params.Sources = pkg.StringPtr(strings.Join(sources, ","))
	}

	payments, metadata, err := s.repo.NonPosted.ListNonPosted(ctx, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]nonPostedResponse, len(payments))

	for idx, p := range payments {
		v, err := s.structureNonPosted(&p, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"data": rsp,
	})
}

func (s *Server) listNonPostedByTransactionSource(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	query := ctx.Param("type")
	if query == "" || (query != "mpesa" && query != "internal") {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid transaction_type")))

		return
	}

	payments, err := s.repo.NonPosted.ListNonPostedByTransactionSource(
		ctx,
		query,
		&pkg.PaginationMetadata{CurrentPage: pageNo},
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]nonPostedResponse, len(payments))

	for idx, p := range payments {
		v, err := s.structureNonPosted(&p, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) listUnassignedNonPostedPayments(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	payments, err := s.repo.NonPosted.ListUnassignedNonPosted(ctx, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]nonPostedResponse, len(payments))

	for idx, p := range payments {
		v, err := s.structureNonPosted(&p, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) getNonPostedPayment(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	payment, err := s.repo.NonPosted.GetNonPosted(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.structureNonPosted(&payment, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

// type assignNonPostedPaymentRequest struct {
// 	ClientID uint32  `binding:"required" json:"client_id"`
// 	AdminID  uint32  `binding:"required" json:"admin_id"`
// 	Amount   float64 `binding:"required" json:"amount"`
// }

// func (s *Server) assignNonPostedPayment(ctx *gin.Context) {
// 	var req assignNonPostedPaymentRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

// 		return
// 	}

// 	id, err := pkg.StringToUint32(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))

// 		return
// 	}

// 	err = s.payments.TriggerManualPayment(ctx, services.ManualPaymentData{
// 		LoanID:      id,
// 		ClientID:    req.ClientID,
// 		Amount:      req.Amount,
// 		AdminUserID: req.AdminID,
// 	})
// 	if err != nil {
// 		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
// }

func (s *Server) structureNonPosted(p *repository.NonPosted, ctx *gin.Context) (nonPostedResponse, error) {
	v := nonPostedResponse{
		ID:                p.ID,
		TransactionSource: string(p.TransactionSource),
		TransactionNumber: p.TransactionNumber,
		AccountNumber:     p.AccountNumber,
		PhoneNumber:       p.PhoneNumber,
		PayingName:        p.PayingName,
		Amount:            p.Amount,
		PaidDate:          p.PaidDate,
	}

	if p.AssignedTo != nil {
		v.Assigned = true
		
		client, err := s.repo.Clients.GetClient(ctx, *p.AssignedTo)
		if err != nil {
			log.Println(*p.AssignedTo)
			return nonPostedResponse{}, err
		}

		branch, err := s.repo.Branches.GetBranchByID(ctx, client.BranchID)
		if err != nil {
			return nonPostedResponse{}, err
		}

		v.AssignedTo = clientShortResponse{
			ID:          client.ID,
			FullName: client.FullName,
			PhoneNumber: client.PhoneNumber,
			Overpayment: client.Overpayment,
			DueAmount: client.DueAmount,
			BranchName: branch.Name,
		}
	}

	return v, nil
}
