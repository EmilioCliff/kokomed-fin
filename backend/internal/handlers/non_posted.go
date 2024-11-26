package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type nonPostedResponse struct {
	ID                uint32            `json:"id"`
	TransactionNumber string            `json:"transaction_number"`
	AccountNumber     string            `json:"account_number"`
	PhoneNumber       string            `json:"phone_number"`
	PayingName        string            `json:"paying_name"`
	Amount            float64           `json:"amount"`
	PaidDate          time.Time         `json:"paid_date"`
	AssignedTo        repository.Client `json:"assigned_to"`
}

func (s *Server) listAllNonPostedPayments(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	payments, err := s.repo.NonPosted.ListNonPosted(ctx, &pkg.PaginationMetadata{CurrentPage: pageNo})
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

type assignNonPostedPaymentRequest struct {
	ClientID uint32 `binding:"required" json:"client_id"`
}

func (s *Server) assignNonPostedPayment(ctx *gin.Context) {
	var req assignNonPostedPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	_, err = s.repo.NonPosted.AssignNonPosted(ctx, id, req.ClientID)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) structureNonPosted(p *repository.NonPosted, ctx *gin.Context) (nonPostedResponse, error) {
	v := nonPostedResponse{
		ID:                p.ID,
		TransactionNumber: p.TransactionNumber,
		AccountNumber:     p.AccountNumber,
		PhoneNumber:       p.PhoneNumber,
		PayingName:        p.PayingName,
		Amount:            p.Amount,
		PaidDate:          p.PaidDate,
	}

	if p.AssignedTo != nil {
		client, err := s.repo.Clients.GetClient(ctx, *p.AssignedTo)
		if err != nil {
			return nonPostedResponse{}, err
		}

		v.AssignedTo = client
	}

	return v, nil
}
