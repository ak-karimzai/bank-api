package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ak-karimzai/ak-karimzai/simpleb/internal/db"
	"github.com/ak-karimzai/ak-karimzai/simpleb/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) createAccount(
	ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errResponse(err))
		return
	}

	authPayload := ctx.MustGet(
		authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	}

	acc, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(
					http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(
	ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errResponse(err))
		return
	}

	acc, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		var status int
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		ctx.JSON(
			status, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(
		authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != acc.Owner {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAccounts(
	ctx *gin.Context) {
	var req getAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errResponse(err))
		return
	}

	authPayload := ctx.MustGet(
		authorizationPayloadKey).(*token.Payload)
	arg := db.GetAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
