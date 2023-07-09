package handlers

import (
	"net/http"
	"time"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/repository"
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tokenMaker  token.Maker
	config      *util.Config
}

func NewUserHandler(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tokenMaker token.Maker,
	config *util.Config,
) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenMaker:  tokenMaker,
		config:      config,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserReponse struct {
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PwdChangedAt time.Time `json:"pwd_changed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func newUserResponse(user db.User) UserReponse {
	return UserReponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Email:        user.Email,
		PwdChangedAt: user.PwdChangedAt,
		CreatedAt:    user.CreatedAt,
	}
}

func (UserHandler *UserHandler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPwd, err := util.HashPasswrod(req.Password)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:  req.Username,
		HashedPwd: hashedPwd,
		FullName:  req.FullName,
		Email:     req.Email,
	}

	user, err := UserHandler.userRepo.CreateUser(ctx, arg)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID   `json:"session_id"`
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	User                  UserReponse `json:"user"`
}

func (userHandler *UserHandler) LoginUser(
	ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := userHandler.userRepo.GetUser(ctx, req.Username)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPwd), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := userHandler.tokenMaker.CreateToken(
		req.Username, userHandler.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	refreshToken, refreshPayload, err := userHandler.tokenMaker.CreateToken(
		req.Username, userHandler.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	session, err := userHandler.sessionRepo.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
