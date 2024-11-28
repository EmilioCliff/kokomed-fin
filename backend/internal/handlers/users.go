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

type userResponse struct {
	ID          uint32    `json:"id"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	BranchName  string    `json:"branch_name"`
	CreatedAt   time.Time `json:"created_at"`
	// RefreshToken string    `json:"refresh_token"`
}

type createUserRequest struct {
	Firstname   string `binding:"required"                   json:"firstname"`
	Lastname    string `binding:"required"                   json:"lastname"`
	PhoneNumber string `binding:"required"                   json:"phone_number"`
	Email       string `binding:"required"                   json:"email"`
	BranchID    uint32 `binding:"required"                   json:"branch_id"`
	Role        string `binding:"required,oneof=admin agent" json:"role"`
	CreatedBy   uint32 `binding:"required"                   json:"created_by"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	password := fmt.Sprintf("%s.%s.%v", req.Firstname, req.Email, req.PhoneNumber[len(req.PhoneNumber)-3:])

	hashPassword, err := pkg.GenerateHashPassword(password, s.config.PASSWORD_COST)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	user, err := s.repo.Users.CreateUser(ctx, &repository.User{
		FullName:     req.Firstname + " " + req.Lastname,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Role:         strings.ToUpper(req.Role),
		BranchID:     req.BranchID,
		CreatedBy:    req.CreatedBy,
		UpdatedAt:    time.Now(),
		UpdatedBy:    req.CreatedBy,
		Password:     hashPassword, // have a default password(firstname.role.last3phonedigits)
		RefreshToken: "defaulted",  // generate his refresh token(first login)
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// TODO: send verification email to user

	v, err := s.convertGeneratedUser(ctx, &user)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

type loginUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type loginUserResponse struct {
	UserData              userResponse `json:"user_data"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	user, err := s.repo.Users.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// create new password(first login)
	if user.PasswordUpdated == 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "password not updated"})

		return
	}

	err = pkg.ComparePasswordAndHash(user.Password, req.Password)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid password")))

		return
	}

	// generate access token and update refresh token
	accesstoken, err := s.maker.CreateToken(user.Email, user.ID, user.BranchID, user.Role, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	refreshToken, err := s.maker.CreateToken(user.Email, user.ID, user.BranchID, user.Role, s.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	_, err = s.repo.Users.UpdateUser(ctx, &repository.UpdateUser{ID: user.ID, RefreshToken: pkg.StringPtr(refreshToken)})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.convertGeneratedUser(ctx, &user)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		UserData:              v,
		AccessToken:           accesstoken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  time.Now().Add(s.config.TOKEN_DURATION),
		RefreshTokenExpiresAt: time.Now().Add(s.config.REFRESH_TOKEN_DURATION),
	})
}

type updateUserCredentialsRequest struct {
	Email       string `binding:"required" json:"email"`
	NewPassword string `binding:"required" json:"new_password"`
}

func (s *Server) updateUserCredentials(ctx *gin.Context) {
	var req updateUserCredentialsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	if payload, ok := ctx.Get("Payload"); ok {
		if payload.(*pkg.Payload).Email != req.Email {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})

			return
		}
	}

	hashpassword, err := pkg.GenerateHashPassword(req.NewPassword, s.config.PASSWORD_COST)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	err = s.repo.Users.UpdateUserPassword(ctx, req.Email, hashpassword)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) refreshToken(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "email is required")))

		return
	}

	user, err := s.repo.Users.GetUserByEmail(ctx, email)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	payload, err := s.maker.GetPayload(user.RefreshToken)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusNotExtended, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "refresh token is expired")))

		return
	}

	accesstoken, err := s.maker.CreateToken(email, user.ID, user.BranchID, user.Role, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.convertGeneratedUser(ctx, &user)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		UserData:              v,
		AccessToken:           accesstoken,
		RefreshToken:          user.RefreshToken,
		AccessTokenExpiresAt:  time.Now().Add(s.config.TOKEN_DURATION),
		RefreshTokenExpiresAt: time.Now().Add(s.config.REFRESH_TOKEN_DURATION),
	})
}

func (s *Server) getUser(ctx *gin.Context) {
	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	user, err := s.repo.Users.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.convertGeneratedUser(ctx, &user)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

func (s *Server) listUsers(ctx *gin.Context) {
	pageNo, err := pkg.StringToUint32(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	users, err := s.repo.Users.ListUsers(ctx, &pkg.PaginationMetadata{CurrentPage: pageNo})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	rsp := make([]userResponse, len(users))

	for idx, u := range users {
		v, err := s.convertGeneratedUser(ctx, &u)
		if err != nil {
			ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

			return
		}

		rsp[idx] = v
	}

	log.Println(rsp)

	ctx.JSON(http.StatusOK, rsp)
}

type updateUserRequest struct {
	Role      string `binding:"required,oneof=ADMIN AGENT" json:"role"`
	BranchID  uint32 `binding:"required"                   json:"branch_id"`
	UpdatedBy uint32 `binding:"required"                   json:"updated_by"`
}

func (s *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	user, err := s.repo.Users.UpdateUser(ctx, &repository.UpdateUser{
		ID:        id,
		Role:      pkg.StringPtr(req.Role),
		BranchID:  pkg.Uint32Ptr(req.BranchID),
		UpdatedBy: pkg.Uint32Ptr(req.UpdatedBy),
		UpdatedAt: pkg.TimePtr(time.Now()),
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	v, err := s.convertGeneratedUser(ctx, &user)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, v)
}

func (s *Server) convertGeneratedUser(ctx *gin.Context, user *repository.User) (userResponse, error) {
	branch, err := s.repo.Branches.GetBranchByID(ctx, user.BranchID)
	if err != nil {
		return userResponse{}, err
	}

	return userResponse{
		ID:          user.ID,
		Fullname:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		BranchName:  branch.Name,
		CreatedAt:   user.CreatedAt,
		// RefreshToken: user.RefreshToken,
	}, nil
}
