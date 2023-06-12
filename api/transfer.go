package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ak-karimzai/ak-karimzai/simpleb/internal/db"
	"github.com/ak-karimzai/ak-karimzai/simpleb/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) createTransfer(
	ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(
	ctx *gin.Context, accID int64, currency string) (db.Account, bool) {
	acc, err := server.store.GetAccount(ctx, accID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return acc, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return acc, false
	}

	if acc.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accID, acc.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return acc, false
	}
	return acc, true
}
