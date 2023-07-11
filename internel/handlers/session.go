package handlers

import (
	"fmt"
	"net/http"
	"time"

	errorhandler "github.com/ak-karimzai/bank-api/internel/error_handler"
	"github.com/ak-karimzai/bank-api/internel/repository"
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionRepo repository.SessionRepository
	tokenMaker  token.Maker
	config      *util.Config
}

func NewSessionHandler(
	sessionRepo repository.SessionRepository,
	tokenMaker token.Maker,
	config *util.Config,
) *SessionHandler {
	return &SessionHandler{
		sessionRepo: sessionRepo,
		tokenMaker:  tokenMaker,
		config:      config,
	}
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (sessionHandler *SessionHandler) RenewAccessToken(
	ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	refreshPayload, err := sessionHandler.
		tokenMaker.
		VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
	}

	session, err := sessionHandler.sessionRepo.
		GetSession(ctx, refreshPayload.ID)
	if err != nil {
		finalErr := errorhandler.DbErrorHandler(err)
		ctx.JSON(toHttpError(finalErr), finalErr)
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(
			http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(
			http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		ctx.JSON(
			http.StatusUnauthorized, errResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(
			http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, accessPayload, err := sessionHandler.tokenMaker.CreateToken(
		refreshPayload.Username, sessionHandler.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
