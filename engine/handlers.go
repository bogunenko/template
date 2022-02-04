package engine

import (
	"database/sql"
	"net/http"

	db "github.com/bogunenko/template/db/sqlc"
	"github.com/gin-gonic/gin"
)

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=0"`
}

func (engine *Engine) createAccount(ctx *gin.Context) {

	res, err := engine.store.CreateAccount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id, _ := res.LastInsertId()
	account, err := engine.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (engine *Engine) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := engine.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type amountRequest struct {
	ID     int64 `json:"id" binding:"required,min=0"`
	Amount int64 `json:"amount" binding:"required,min=0"`
}

func (engine *Engine) deposit(ctx *gin.Context) {
	var req amountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := engine.store.Deposit(ctx, db.DepositParams{ID: req.ID, Amount: req.Amount})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account, err := engine.store.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (engine *Engine) withdraw(ctx *gin.Context) {
	var req amountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := engine.store.Withdraw(ctx, db.WithdrawParams{ID: req.ID, Amount: req.Amount})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account, err := engine.store.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type transferRequest struct {
	FromID int64 `json:"from_id" binding:"required,min=0"`
	ToID   int64 `json:"to_id" binding:"required,min=0"`
	Amount int64 `json:"amount" binding:"required,min=0"`
}

func (engine *Engine) transfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := engine.store.Transfer(ctx, db.CreateTransactionParams{FromID: req.FromID, ToID: req.ToID, Amount: req.Amount})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, _ := result.LastInsertId()
	transaction, err := engine.store.GetTransaction(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}
