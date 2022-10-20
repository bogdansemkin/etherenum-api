package http

import (
	"etherenum-api/etherenum-service/api/internal/entities"
	"etherenum-api/etherenum-service/api/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type blockChainController struct {
	Controller
}

func NewBlockchainRoutes(options ControllerOptions) {
	r := &blockChainController{Controller{
		Config:  options.Config,
		Service: options.Service,
		Repos:   options.Repos,
		Logger:  options.Logger.Named("blockchainRoutes"),
	}}

	p := options.Handler.Group("/transactions")
	{
		p.GET("/", errorHandler(options, r.getTransactions))
		p.GET("/:filter", errorHandler(options, r.getTransactionByFilter))
	}
}

type getTransactionsQuery struct {
	Page int64 `form:"page" json:"page"`
}

type getTransactionsResponse struct {
	Transactions *[]entities.Transaction
}

func (b *blockChainController) getTransactions(c *gin.Context) (interface{}, *Error) {
	logger := b.Logger.Named("getTransactions").WithContext(c)

	var query getTransactionsQuery
	err := c.BindQuery(&query)
	if err != nil {
		logger.Info("error during binding query", "err", err)
		return nil, &Error{Type: ErrorTypeClient, Message: "error during binding query", Details: err}
	}
	logger.With("query", "query", query.Page)
	logger.Debug("query", "query", query.Page)

	transactions, err := b.Controller.Service.Transaction.GetAll(c, query.Page)
	if err != nil {
		switch err.(type) {
		//case service.LessPointerZeroError:
		//	fmt.Printf("less pointer zero exception, %s\n", err)
		//	return
		default:
			logger.Info("error on get all transaction", "err", err)
			return nil, &Error{Type: ErrorTypeServer, Message: "error during getting all transactions", Details: err}
		}
	}

	fmt.Println("Successfully got all transactions")
	return getTransactionsResponse{Transactions: transactions}, nil
}

type getTransactionByFilterQuery struct {
	Page int64 `form:"page" json:"page"`
}

type getTransactionByFilterParams struct {
	Filter string `uri:"filter" json:"filter" binding:"required"`
}

type getTransactionByFilterResponse struct {
	Trans []entities.Transaction
}

func (b *blockChainController) getTransactionByFilter(c *gin.Context) (interface{}, *Error) {
	logger := b.Logger.Named("getTransactionByFilter").WithContext(c)

	var query getTransactionByFilterQuery
	var body getTransactionByFilterParams

	err := c.BindQuery(&query)
	if err != nil {
		logger.Info("error during binding query", "err", err)
		return nil, &Error{Type: ErrorTypeClient, Message: "error during binding query", Details: err}
	}

	err = c.ShouldBindUri(&body)
	if err != nil {
		logger.Info("error on get transaction by filter", "err", err)
		return nil, &Error{Type: ErrorTypeClient, Message: "error on get transactions by filter", Details: err}
	}
	logger.With("body", "body", body)
	logger.Debug("body", "body", body)

	transactions, err := b.Service.Transaction.GetByFilter(c, body.Filter, query.Page)
	if err != nil {
		switch err.(type) {
		case service.NilPointerDataError:
			logger.Info("nil pointer data error", "err", err)
			return nil, &Error{Type: ErrorTypeClient, Message: "wrong filter for transaction", Details: err}
		default:
			logger.Info("error during get data by filter", "err", err)
			return nil, &Error{Type: ErrorTypeServer, Message: "error during getting transaction", Details: err}
		}
	}

	logger.Info("Successfully got transaction(s) by filter")
	return getTransactionByFilterResponse{Trans: transactions.Trans}, nil
}
