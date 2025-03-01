package handlers

import (
	"fmt"
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
	log.Println("cache miss")
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

	params := repository.NonPostedCategory{}
	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
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

	response := gin.H{
		"metadata": metadata,
		"data": rsp,
	}

	cacheKey := constructCacheKey("non-posted/all", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
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

type listClientsNonPostedReq struct {
	ID uint32 `json:"id"`
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

	rslt, pgData, err := s.repo.NonPosted.GetClientNonPosted(ctx, req.ID, req.PhoneNumber, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data": rslt,
	})
}

func (s *Server) structureNonPosted(p *repository.NonPosted, ctx *gin.Context) (nonPostedResponse, error) {
	cacheKey := fmt.Sprintf("non-posted:%v", p.ID)
	var dataCached nonPostedResponse

	exists, _ := s.cache.Get(ctx, cacheKey, &dataCached)
	if exists {
		return dataCached, nil
	}
	
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

	if err := s.cache.Set(ctx, cacheKey, v, 3*time.Minute); err != nil {
		return nonPostedResponse{}, err
	}

	return v, nil
}
