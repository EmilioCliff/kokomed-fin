package handlers

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type clientResponse struct {
	ID            uint32            `json:"id"`
	FullName          string            `json:"fullName"`
	PhoneNumber   string            `json:"phoneNumber"`
	IdNumber      string            `json:"idNumber"`
	Dob           string            `json:"dob"`
	Gender        string            `json:"gender"`
	Active        bool              `json:"active"`
	BranchName    string            `json:"branchName"`
	AssignedStaff userShortResponse `json:"assignedStaff"` 
	Overpayment   float64           `json:"overpayment"`
	DueAmount float64 `json:"dueAmount"`
	CreatedBy     userShortResponse `json:"createdBy"` 
	CreatedAt     time.Time         `json:"createdAt"`
}

type userShortResponse struct {
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type clientShortResponse struct {
	Overpayment   float64           `json:"overpayment"`
	DueAmount float64 `json:"dueAmount"`
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	BranchName string 	`json:"branchName"`
}

type createClientRequest struct {
	FirstName     string `binding:"required"                   json:"firstName"`
	LastName      string `binding:"required"                   json:"lastName"`
	PhoneNumber   string `binding:"required"                   json:"phoneNumber"`
	IdNumber      string `                                     json:"idNumber"`
	Dob           string `                                     json:"dob"`
	Gender        string `binding:"required,oneof=MALE FEMALE" json:"gender"`
	Active        bool   `                                     json:"active"`
	BranchID      uint32 `binding:"required"                   json:"branchId"`
	AssignedStaff uint32 `binding:"required"                   json:"assignedStaffId"`
}

func (s *Server) createClient(ctx *gin.Context) {
	var req createClientRequest
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

	params := &repository.Client{
		FullName:      req.FirstName + " " + req.LastName,
		PhoneNumber:   req.PhoneNumber,
		Gender:        req.Gender,
		BranchID:      req.BranchID,
		AssignedStaff: req.AssignedStaff,
		UpdatedBy:     payloadData.UserID,
		CreatedBy:     payloadData.UserID,
	}

	if req.Dob != "" {
		dob, err := time.Parse("2006-01-02", req.Dob)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.Dob = pkg.TimePtr(dob)
	}

	if req.IdNumber != "" {
		params.IdNumber = pkg.StringPtr(req.IdNumber)
	}

	client, err := s.repo.Clients.CreateClient(ctx, params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.structureClient(&client, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

type updateClient struct {
	IdNumber      string `json:"idNumber"`
	Dob           string `json:"dob"`
	Active        string   `json:"active"`
	BranchID      uint32 `json:"branchId"`
}

func (s *Server) updateClient(ctx *gin.Context) {
	var req updateClient
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
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

	params := &repository.UpdateClient{
		ID: id,
		UpdatedBy: payloadData.UserID,
	}

	if req.Dob != "" {
		dob, err := time.Parse("2006-01-02", req.Dob)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

			return
		}

		params.Dob = pkg.TimePtr(dob)
	}

	if req.IdNumber != "" {
		params.IdNumber = pkg.StringPtr(req.IdNumber)
	}

	if req.Active != "" {
		if req.Active == "true" {
			params.Active = pkg.BoolPtr(true)
		} else {
			params.Active = pkg.BoolPtr(false)
		}
	}

	if req.BranchID != 0 {
		params.BranchID = pkg.Uint32Ptr(req.BranchID)
	}

	err = s.repo.Clients.UpdateClient(ctx, params)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Client Updated"})
}

func (s *Server) listClients(ctx *gin.Context) {
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

	params := repository.ClientCategorySearch{}

	name := ctx.Query("search")
	if name != "" {
		params.Search = pkg.StringPtr(name)
	}

	active := ctx.Query("active")
	if active != "" {
		if active == "2" {
			params.Active = pkg.BoolPtr(false)
		} else {
			params.Active = pkg.BoolPtr(true)
		}
	}

	clients, metadata, err := s.repo.Clients.ListClients(ctx, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]clientResponse, len(clients))

	for idx, c := range clients {
		v, err := s.structureClient(&c, ctx)
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

func (s *Server) getClient(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	client, err := s.repo.Clients.GetClient(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.structureClient(&client, ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

func (s *Server) listClientsByBranch(ctx *gin.Context) {
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

	clients, err := s.repo.Clients.ListClientsByBranch(ctx, id, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]clientResponse, len(clients))

	for idx, c := range clients {
		v, err := s.structureClient(&c, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) listClientsByActive(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	check := true

	status := ctx.DefaultQuery("status", "active")
	if status == "inactive" {
		check = false
	}

	clients, err := s.repo.Clients.ListClientsByActiveStatus(ctx, check, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]clientResponse, len(clients))

	for idx, c := range clients {
		v, err := s.structureClient(&c, ctx)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *Server) structureClient(c *repository.Client, ctx *gin.Context) (clientResponse, error) {
	assignedStaff, err := s.repo.Users.GetUserByID(ctx, c.AssignedStaff)
	if err != nil {
		return clientResponse{}, err
	}

	// updatedBy, err := s.repo.Users.GetUserByID(ctx, c.UpdatedBy)
	// if err != nil {
	// 	return clientResponse{}, err
	// }

	createdBy, err := s.repo.Users.GetUserByID(ctx, c.CreatedBy)
	if err != nil {
		return clientResponse{}, err
	}

	branch, err := s.repo.Branches.GetBranchByID(ctx, c.BranchID)
	if err != nil {
		return clientResponse{}, err
	}

	rsp := clientResponse{
		ID:            c.ID,
		FullName:          c.FullName,
		PhoneNumber:   c.PhoneNumber,
		Gender:        c.Gender,
		Active:        c.Active,
		BranchName:    branch.Name,
		AssignedStaff: userShortResponse{ID: assignedStaff.ID, FullName: assignedStaff.FullName, Email: assignedStaff.Email, PhoneNumber: assignedStaff.PhoneNumber},
		Overpayment:   c.Overpayment,
		CreatedBy:     userShortResponse{ID: createdBy.ID, FullName: createdBy.FullName, Email: createdBy.Email, PhoneNumber: createdBy.PhoneNumber},
		CreatedAt:     c.CreatedAt,
		DueAmount: c.DueAmount,
	}
	// UpdatedAt:     c.UpdatedAt,
	// UpdatedBy:     userShortResponse{ID: updatedBy.ID, FullName: updatedBy.FullName},

	if c.Dob != nil {
		rsp.Dob = c.Dob.Format("2006-01-02")
	}

	if c.IdNumber != nil {
		rsp.IdNumber = *c.IdNumber
	}

	return rsp, nil
}
