package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
)

type userResponse struct {
	ID          uint32    `json:"id"`
	Fullname    string    `json:"fullName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Role        string    `json:"role"`
	BranchName  string    `json:"branchName"`
	CreatedAt   time.Time `json:"createdAt"`
}

type createUserRequest struct {
	Firstname   string `binding:"required"                   json:"firstName"`
	Lastname    string `binding:"required"                   json:"lastName"`
	PhoneNumber string `binding:"required"                   json:"phoneNumber"`
	Email       string `binding:"required"                   json:"email"`
	BranchID    uint32 `binding:"required"                   json:"branchId"`
	Role        string `binding:"required,oneof=ADMIN AGENT" json:"role"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	password := fmt.Sprintf("%s.%s.%v", req.Firstname, req.Role, req.PhoneNumber[len(req.PhoneNumber)-3:])
	log.Print(password)

	hashPassword, err := pkg.GenerateHashPassword(password, s.config.PASSWORD_COST)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

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

	user, err := s.repo.Users.CreateUser(ctx, &repository.User{
		FullName:     req.Firstname + " " + req.Lastname,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Role:         strings.ToUpper(req.Role),
		BranchID:     req.BranchID,
		CreatedBy:    payloadData.UserID,
		UpdatedAt:    time.Now(),
		UpdatedBy:    payloadData.UserID,
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
	UserData              userResponse `json:"userData"`
	AccessToken           string       `json:"accessToken"`
	RefreshToken          string       `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time    `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time    `json:"refreshTokenExpiresAt"`
}

type forgotPasswordRequest struct {
	Email string `binding:"required" json:"email"`
}

func (s *Server) forgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	// check if user exists in db
	exists := s.repo.Users.CheckUserExistance(ctx, req.Email)
	if !exists {
		ctx.JSON(http.StatusNotFound, errorResponse(pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")))

		return
	}

	err := s.worker.DistributeTaskSendResetPassword(ctx, services.SendResetPasswordPayload{
		Email: req.Email,
	})
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "email sent"})
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

	ctx.SetCookie("refreshToken", refreshToken, int(s.config.REFRESH_TOKEN_DURATION), "/", "", true, true)

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

func (s *Server) logoutUser(ctx *gin.Context) {
	ctx.SetCookie("refreshToken", "", -1, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

type updateUserCredentialsRequest struct {
	NewPassword string `binding:"required" json:"newPassword"`
}

func (s *Server) updateUserCredentials(ctx *gin.Context) {
	var req updateUserCredentialsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	token := ctx.Param("token")
	payload, err := s.maker.VerifyToken(token)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	// get token from redis db
	// send a uuid then connect uuid to the redis key value
	// ttl is 10 min expiration of the token

	hashpassword, err := pkg.GenerateHashPassword(req.NewPassword, s.config.PASSWORD_COST)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	err = s.repo.Users.UpdateUserPassword(ctx, payload.Email, hashpassword)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) refreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token not found"})
		return
	}

	payload, err := s.maker.GetPayload(refreshToken)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	if payload.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		ctx.JSON(http.StatusNotExtended, errorResponse(pkg.Errorf(pkg.INVALID_ERROR, "refresh token is expired")))

		return
	}

	accesstoken, err := s.maker.CreateToken(payload.Email, payload.UserID, payload.BranchID, payload.Role, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		AccessToken:           accesstoken,
		AccessTokenExpiresAt:  time.Now().Add(s.config.TOKEN_DURATION),
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

	params := repository.CategorySearch{}
	cacheParams := map[string][]string{
		"page": {pageNoStr},
		"limit": {pageSizeStr},
	}

	name := ctx.Query("search")
	if name != "" {
		params.Search = pkg.StringPtr(name)
		cacheParams["search"] = []string{name}
	}

	role := ctx.Query("role")
	if role != "" {
		roles := strings.Split(role, ",")

		for i := range roles {
			roles[i] = strings.TrimSpace(roles[i])
		}

		params.Role = pkg.StringPtr(strings.Join(roles, ","))
		cacheParams["role"] = []string{strings.Join(roles, ",")}
	}

	users, metadata, err := s.repo.Users.ListUsers(ctx, &params, &pkg.PaginationMetadata{CurrentPage: pageNo, PageSize: pageSize})
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

	response := gin.H{
		"metadata": metadata,
		"data": rsp,
	}

	cacheKey := constructCacheKey("user", cacheParams)

	err = s.cache.Set(ctx, cacheKey, response, 20*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.Errorf(pkg.INTERNAL_ERROR, "failed caching: %s", err))

		return
	}

	ctx.JSON(http.StatusOK, response)
}

type updateUserRequest struct {
	Role      string `json:"role"`
	BranchID  uint32 `json:"branchId"`
}

func (s *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
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

	id, err := pkg.StringToUint32(ctx.Param("id"))
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(pkg.Errorf(pkg.INVALID_ERROR, err.Error())))

		return
	}

	params := repository.UpdateUser{
		ID: id,
		UpdatedBy: pkg.Uint32Ptr(payloadData.UserID),
		UpdatedAt: pkg.TimePtr(time.Now()),
	}

	if req.Role != "" {
		params.Role = pkg.StringPtr(req.Role)
	}

	if req.BranchID != 0 {
		params.BranchID = pkg.Uint32Ptr(req.BranchID)
	}

	user, err := s.repo.Users.UpdateUser(ctx, &params)
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
	cacheKey := fmt.Sprintf("user:%v", user.ID)
	var dataCached userResponse

	exists, _ := s.cache.Get(ctx, cacheKey, &dataCached)
	if exists {
		log.Println("Cached Hit: ", cacheKey)
		return dataCached, nil
	}

	branch, err := s.repo.Branches.GetBranchByID(ctx, user.BranchID)
	if err != nil {
		return userResponse{}, err
	}

	rsp := userResponse{
		ID:          user.ID,
		Fullname:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		BranchName:  branch.Name,
		CreatedAt:   user.CreatedAt,
		// RefreshToken: user.RefreshToken,
	}

	if err := s.cache.Set(ctx, cacheKey, rsp, 3*time.Minute); err != nil {
		return userResponse{}, err
	}

	return rsp, nil
}
