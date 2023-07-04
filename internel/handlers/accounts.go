package handlers

import (
	"errors"
	"net/http"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/middlewares"
	"github.com/ak-karimzai/bank-api/internel/repository"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountRepo repository.AccountRepository
}

func NewAccountHandler(accountRepo repository.AccountRepository) *AccountHandler {
	return &AccountHandler{
		accountRepo: accountRepo,
	}
}

type CreateAccountReq struct {
	Owner    string      `json:"owner" binding:"required"`
	Currency db.Currency `json:"currency" binding:"required"`
}

func (accHandler *AccountHandler) CreateAccount(ctx *gin.Context) {
	var req CreateAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}
	acc, err := accHandler.accountRepo.CreateAccount(ctx, arg)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,numeric,gt=0"`
}

func (accHandler *AccountHandler) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := middlewares.GetPayload(ctx)
	acc, err := accHandler.accountRepo.GetAccount(ctx, req.ID)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}

	if acc.Owner != payload.Username {
		err := errors.New("account doesn't  belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type listAccountsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (accHandler *AccountHandler) ListAccounts(ctx *gin.Context) {
	var req listAccountsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload := middlewares.GetPayload(ctx)
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner:  payload.Username,
	}

	acc, err := accHandler.accountRepo.ListAccounts(ctx, arg)
	if err != nil {
		finalErr := dbErrorHandler(err)
		ctx.JSON(finalErr.Status, finalErr)
		return
	}
	ctx.JSON(http.StatusOK, acc)
}
