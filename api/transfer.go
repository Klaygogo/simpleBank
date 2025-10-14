package api

import (
	"fmt"
	db "github.com/Klaygogo/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func (server *Server) validAccount(c *gin.Context, accountID int64, currency string) bool {
	accout, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if accout.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: got=%s, want=%s", accountID, accout.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// from account
	if !server.validAccount(c, req.FromAccountID, req.Currency) {
		return
	}
	// to account
	if !server.validAccount(c, req.ToAccountID, req.Currency) {
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	if req.FromAccountID == req.ToAccountID {
		err := fmt.Errorf("cannot transfer to the same account")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, _ := server.store.GetAccount(c, req.FromAccountID)
	if arg.Amount > fromAccount.Balance {
		err := fmt.Errorf("insufficient balance")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	c.JSON(http.StatusOK, result)
}

type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(c *gin.Context) {
	var req getTransferRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(c, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))

	}
	c.JSON(http.StatusOK, transfer)
}

type listTransferRequest struct {
	PageID        int32 `form:"page_id" binding:"required,min=1"`
	PageSize      int32 `form:"page_size" binding:"required,min=5,max=10"`
	FromAccountID int64 `form:"from_account_id" binding:"required"`
	ToAccountID   int64 `form:"to_account_id" binding:"required"`
}

func (server *Server) listTransfer(c *gin.Context) {
	var req listTransferRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.ListTransfersParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}
	transfers, err := server.store.ListTransfers(c, args)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transfers)
}
