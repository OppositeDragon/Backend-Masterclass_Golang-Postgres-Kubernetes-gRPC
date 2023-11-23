package api

import (
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"fromAccountId" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountId" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	isFromAccountValid  := server.isValidAccountForCurrency(ctx, req.FromAccountID, req.Currency) 
	isToAccountValid := server.isValidAccountForCurrency(ctx, req.ToAccountID, req.Currency)

	if isFromAccountValid!=nil||nil!= isToAccountValid {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("%s, %s",isFromAccountValid,isToAccountValid)))
		return 	
	} 

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) isValidAccountForCurrency(ctx *gin.Context, accountId int64, currency string)  error {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		return   err
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		return  err
	}
	return  nil
}
