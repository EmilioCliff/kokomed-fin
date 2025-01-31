package handlers

import (
	"net/http"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type productResponse struct {
	ID             uint32    `json:"id"`
	BranchName     string    `json:"branchName"`
	LoanAmount     float64   `json:"loanAmount"`
	RepayAmount    float64   `json:"repayAmount"`
	InterestAmount float64   `json:"interestAmount"`
}
// UpdatedBy      uint32    `json:"updated_by"`
// UpdatedAt      time.Time `json:"updated_at"`
// CreatedAt      time.Time `json:"created_at"`

type createProductRequest struct {
	BranchID    uint32  `binding:"required" json:"branchId"`
	LoanAmount  float64 `binding:"required" json:"loanAmount"`
	RepayAmount float64 `binding:"required" json:"repayAmount"`
}

func (s *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

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

	product, err := s.repo.Products.CreateProduct(ctx, &repository.Product{
		BranchID:       req.BranchID,
		LoanAmount:     req.LoanAmount,
		RepayAmount:    req.RepayAmount,
		UpdatedBy:      payloadData.UserID,
		InterestAmount: req.RepayAmount - req.LoanAmount,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.structureProduct(&product, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

func (s *Server) getProduct(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	product, err := s.repo.Products.GetProductByID(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.structureProduct(&product, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

func (s *Server) listProducts(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	pageSize, err := pkg.StringToUint32(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	products, pgData, err := s.repo.Products.GetAllProducts(ctx, pkg.StringPtr(ctx.Query("search")), &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]productResponse, len(products))

	for idx, p := range products {
		v, err := s.structureProduct(&p, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, gin.H{
		"metadata": pgData,
		"data": rsp,
	})
}

func (s *Server) listProductsByBranch(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	products, err := s.repo.Products.ListProductByBranch(ctx, id, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]productResponse, len(products))

	for idx, p := range products {
		v, err := s.structureProduct(&p, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) structureProduct(p *repository.Product, ctx *gin.Context) (productResponse, error) {
	branch, err := s.repo.Branches.GetBranchByID(ctx, p.BranchID)
	if err != nil {
		return productResponse{}, err
	}

	return productResponse{
		ID:             p.ID,
		LoanAmount:     p.LoanAmount,
		BranchName:     branch.Name,
		RepayAmount:    p.RepayAmount,
		InterestAmount: p.InterestAmount,
		}, nil
		// UpdatedBy:      p.UpdatedBy,
		// UpdatedAt:      p.UpdatedAt,
		// CreatedAt:      p.CreatedAt,
}
