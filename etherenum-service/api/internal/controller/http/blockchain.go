package http

import (
	"etherenum-api/etherenum-service/api/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type blockChainController struct {
	Controller
}

func NewBlockchainRoutes(options ControllerOptions) {
	r := &blockChainController{Controller{
		Config:    options.Config,
		Service:   options.Service,
		Repos:     options.Repos,
		Logger:    options.Logger.Named("blockchainRoutes"),
		Converter: options.Converter,
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

func (b *blockChainController) getTransactions(c *gin.Context) (interface{}, *Error) {
	logger := b.Logger.Named("getTransactions").WithContext(c)

	var query getTransactionsQuery
	var response []TransactionResponse

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
	for i := range transactions.Trans {
		response = append(response, TransactionResponse{
			Hash:        transactions.Trans[i].Hash,
			From:        transactions.Trans[i].From,
			To:          transactions.Trans[i].To,
			BlockNumber: b.Converter.StringToInt(transactions.Trans[i].BlockNumber),
			Commission:  transactions.Trans[i].Commission,
			Value:       transactions.Trans[i].Value,
			Timestamp:   transactions.Trans[i].Timestamp,
		})
	}

	fmt.Println("Successfully got all transactions")
	return response, nil
}

type getTransactionByFilterQuery struct {
	Page int64 `form:"page" json:"page"`
}

type getTransactionByFilterParams struct {
	Filter string `uri:"filter" json:"filter" binding:"required"`
}

func (b *blockChainController) getTransactionByFilter(c *gin.Context) (interface{}, *Error) {
	logger := b.Logger.Named("getTransactionByFilter").WithContext(c)

	var query getTransactionByFilterQuery
	var body getTransactionByFilterParams
	var response []TransactionResponse

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
	for i := range transactions.Trans {
		response = append(response, TransactionResponse{
			Hash:        transactions.Trans[i].Hash,
			From:        transactions.Trans[i].From,
			To:          transactions.Trans[i].To,
			BlockNumber: b.Converter.StringToInt(transactions.Trans[i].BlockNumber),
			Commission:  transactions.Trans[i].Commission,
			Value:       transactions.Trans[i].Value,
			Timestamp:   transactions.Trans[i].Timestamp,
		})
	}
	logger.Info("Successfully got transaction(s) by filter")
	return response, nil
}
