package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type nonPostedResponse struct {
	ID                uint32              `json:"id"`
	TransactionSource string              `json:"transactionSource"`
	TransactionNumber string              `json:"transactionNumber"`
	AccountNumber     string              `json:"accountNumber"`
	PhoneNumber       string              `json:"phoneNumber"`
	PayingName        string              `json:"payingName"`
	Amount            float64             `json:"amount"`
	PaidDate          time.Time           `json:"paidDate"`
	AssignedTo        clientShortResponse `json:"assignedTo"`
	Assigned          bool                `json:"assigned"`
	AssignedBy        string              `json:"assignedBy"`
}

func (s *Server) listAllNonPostedPayments(ctx *gin.Context) {
	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	fromDateStr := ctx.DefaultQuery("from", "01/01/2025")
	fromDate, err := time.Parse("01/02/2006", fromDateStr)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(
				pkg.Errorf(pkg.INTERNAL_ERROR, "error parsing from date: %s", err.Error()),
			),
		)

		return
	}

	toDateStr := ctx.DefaultQuery("to", time.Now().Format("01/02/2006"))
	toDate, err := time.Parse("01/02/2006", toDateStr)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(
				pkg.Errorf(pkg.INTERNAL_ERROR, "error parsing from date: %s", err.Error()),
			),
		)

		return
	}
	toDate = toDate.Add(24 * time.Hour)

	params := repository.NonPostedCategory{}
	cacheParams := map[string][]string{
		"page":     {pageNoStr},
		"limit":    {pageSizeStr},
		"fromDate": {fromDateStr},
		"toDate":   {toDateStr},
	}

	search := ctx.Query("search")
	if search != "" {
		params.Search = pkg.StringPtr(strings.ToLower(search))
		cacheParams["search"] = []string{search}
	}

	source := ctx.Query("source")
	if source != "" {
		sources := strings.Split(source, ",")

		for i := range sources {
			sources[i] = strings.TrimSpace(sources[i])
		}

		params.Sources = pkg.StringPtr(strings.Join(sources, ","))
		cacheParams["source"] = []string{strings.Join(sources, ",")}
	}

	payments, metadata, err := s.repo.NonPosted.ListNonPosted(
		ctx,
		&params,
		&pkg.PaginationMetadata{
			CurrentPage: pageNo,
			PageSize:    pageSize,
			FromDate:    &fromDate,
			ToDate:      &toDate,
		},
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]nonPostedResponse, len(payments))

	for idx, p := range payments {
		rsp[idx] = structureNonPosted(&p)
	}

	response := gin.H{
		"metadata": metadata,
		"data":     rsp,
	}

	cacheKey := constructCacheKey("non-posted/all", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err),
		)

		return
	}

	ctx.JSON(http.StatusOK, response)
}

type listClientsNonPostedReq struct {
	ID          uint32 `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
}

func (s *Server) listClientsNonPosted(ctx *gin.Context) {
	var req listClientsNonPostedReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageNoStr := ctx.DefaultQuery("page", "1")
	pageNo, err := pkg.StringToUint32(pageNoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	pageSizeStr := ctx.DefaultQuery("limit", "10")
	pageSize, err := pkg.StringToUint32(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	rslt, pgData, err := s.repo.NonPosted.GetClientNonPosted(
		ctx,
		req.ID,
		req.PhoneNumber,
		&pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize},
	)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data":     rslt,
	})
}

func (s *Server) getNonPosted(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	nonPosted, err := s.repo.NonPosted.GetNonPosted(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	response := gin.H{
		"data": nonPosted,
	}

	err = s.cache.Set(
		ctx,
		constructCacheKey(fmt.Sprintf("non-posted/%d", id), map[string][]string{}),
		response,
		1*time.Minute,
	)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err),
		)

		return
	}

	ctx.JSON(http.StatusOK, response)
}

func structureNonPosted(p *repository.NonPosted) nonPostedResponse {
	rsp := nonPostedResponse{
		ID:                p.ID,
		TransactionSource: string(p.TransactionSource),
		TransactionNumber: p.TransactionNumber,
		AccountNumber:     p.AccountNumber,
		PhoneNumber:       p.PhoneNumber,
		PayingName:        p.PayingName,
		Amount:            p.Amount,
		PaidDate:          p.PaidDate,
		AssignedBy:        p.AssignedBy,
	}

	if p.AssignedTo != nil && p.AssignedClient.ID > 0 {
		rsp.Assigned = true
		rsp.AssignedTo = clientShortResponse{
			ID:          p.AssignedClient.ID,
			FullName:    p.AssignedClient.FullName,
			PhoneNumber: p.AssignedClient.PhoneNumber,
			Overpayment: p.AssignedClient.Overpayment,
			BranchName:  p.AssignedClient.BranchName,
		}
	}

	return rsp
}
