package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type createBranchRequest struct {
	Name string `binding:"required" json:"name"`
}

func (s *Server) createBranch(ctx *gin.Context) {
	var req createBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	tc, span := s.tracer.Start(ctx.Request.Context(), "Creating Branch", oteltrace.WithAttributes(attribute.String("name", req.Name)))
	defer span.End()

	branch, err := s.repo.Branches.CreateBranch(tc, &repository.Branch{Name: req.Name})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	s.cache.DelAll(ctx, "branch:limit=*")

	ctx.JSON(http.StatusOK, branch)
}

func (s *Server) getBranch(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	branch, err := s.repo.Branches.GetBranchByID(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, branch)
}

func (s *Server) listBranches(ctx *gin.Context) {
	tc, span := s.tracer.Start(ctx.Request.Context(), "Listing Branches")
	defer span.End()

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

	span.SetAttributes(
		attribute.String("page_no", pageNoStr),
		attribute.String("page_size", pageSizeStr),
	)

	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
	}

	search := ctx.Query("search")
	if search != "" {
		span.SetAttributes(attribute.String("searched", search))
		cacheParams["search"] = []string{search}
	}

	branches, metadata, err := s.repo.Branches.ListBranches(tc, pkg.StringPtr(search), &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	response := gin.H{
		"metadata": metadata,
		"data": branches,
	}

	cacheKey := constructCacheKey("branch", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 1*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) updateBranch(ctx *gin.Context) {
	var req createBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	tc, span := s.tracer.Start(ctx.Request.Context(), "Updating Branch Name", oteltrace.WithAttributes(attribute.String("new_branchName", req.Name), attribute.Int64("branch_id", int64(id))))
	defer span.End()

	branch, err := s.repo.Branches.UpdateBranch(tc, req.Name, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, branch)
}
