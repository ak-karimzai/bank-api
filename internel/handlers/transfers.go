package handlers

import (
	"net/http"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/middlewares"
	"github.com/ak-karimzai/bank-api/internel/repository"
	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferRepo repository.TransferRepository
}

func NewTransferHandler(transferRepo repository.TransferRepository) *TransferHandler {
	return &TransferHandler{
		transferRepo: transferRepo,
	}
}

type CreateTransferReq struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

func (transferHandler *TransferHandler) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := middlewares.GetPayload(ctx)
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Username:      payload.Username,
	}
	result, err := transferHandler.transferRepo.TransferTx(ctx, arg)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
