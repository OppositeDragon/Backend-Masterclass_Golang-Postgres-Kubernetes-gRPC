package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"

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
	fromAccount, errFrom  := server.isValidAccountForCurrency(ctx, req.FromAccountID, req.Currency) 
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Username!=authPayload.Username {
		err:=errors.New("from account does not belong by the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("%s, %s",errFrom,err)))
		return 
	}
	_,errTo := server.isValidAccountForCurrency(ctx, req.ToAccountID, req.Currency)
	if errFrom!=nil||nil!= errTo {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("%s, %s",errFrom,errTo)))
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

func (server *Server) isValidAccountForCurrency(ctx *gin.Context, accountId int64, currency string)  (db.Account,error) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		return  account, err
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		return  account,err
	}
	return  account,nil
}
