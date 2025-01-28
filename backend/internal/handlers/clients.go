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
	FullName          string            `json:"full_name"`
	PhoneNumber   string            `json:"phone_number"`
	IdNumber      string            `json:"id_number"`
	Dob           string            `json:"dob"`
	Gender        string            `json:"gender"`
	Active        bool              `json:"active"`
	BranchName    string            `json:"branch_name"`
	AssignedStaff userShortResponse `json:"assigned_staff"` 
	Overpayment   float64           `json:"overpayment"`
	DueAmount float64 `json:"due_amount"`
	CreatedBy     userShortResponse `json:"created_by"` 
	CreatedAt     time.Time         `json:"created_at"`
}
// UpdatedBy     userShortResponse `json:"updated_by"` 
// UpdatedAt     time.Time         `json:"updated_at"`

type userShortResponse struct {
	ID       uint32 `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type createClientRequest struct {
	FirstName     string `binding:"required"                   json:"first_name"`
	LastName      string `binding:"required"                   json:"last_name"`
	PhoneNumber   string `binding:"required"                   json:"phone_number"`
	IdNumber      string `                                     json:"id_number"`
	Dob           string `                                     json:"dob"`
	Gender        string `binding:"required,oneof=MALE FEMALE" json:"gender"`
	Active        bool   `                                     json:"active"`
	BranchID      uint32 `binding:"required"                   json:"branch_id"`
	AssignedStaff uint32 `binding:"required"                   json:"assigned_staff"`
	// UpdatedBy     uint32 `binding:"required"                   json:"updated_by"`
}

func (s *Server) createClient(ctx *gin.Context) {
	var req createClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	payloadString, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusForbidden, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "No payload found in context")))

		return
	}

	payload, _ := payloadString.(pkg.Payload)

	params := &repository.Client{
		FullName:      req.FirstName + " " + req.LastName,
		PhoneNumber:   req.PhoneNumber,
		Gender:        req.Gender,
		Active:        req.Active,
		BranchID:      req.BranchID,
		AssignedStaff: req.AssignedStaff,
		UpdatedBy:     payload.UserID,
		CreatedBy:     payload.UserID,
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

func (s *Server) updateClient(ctx *gin.Context) {
	var req createClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	payloadString, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.JSON(http.StatusForbidden, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "No payload found in context")))

		return
	}

	payload, _ := payloadString.(pkg.Payload)

	params := &repository.Client{
		ID:            id,
		FullName:      req.FirstName + " " + req.LastName,
		PhoneNumber:   req.PhoneNumber,
		Gender:        req.Gender,
		Active:        req.Active,
		BranchID:      req.BranchID,
		AssignedStaff: req.AssignedStaff,
		UpdatedBy:     payload.UserID,
		CreatedBy:     payload.UserID,
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

	client, err := s.repo.Clients.UpdateClient(ctx, params)
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

func (s *Server) listClients(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	clients, err := s.repo.Clients.ListClients(ctx, &pkg.PaginationMetadata{CurrentPage: pageNo})
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
		AssignedStaff: userShortResponse{ID: assignedStaff.ID, FullName: assignedStaff.FullName},
		Overpayment:   c.Overpayment,
		CreatedBy:     userShortResponse{ID: createdBy.ID, FullName: createdBy.FullName},
		CreatedAt:     c.CreatedAt,
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
