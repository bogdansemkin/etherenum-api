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
	var query getTransactionsQuery
	err := c.BindQuery(&query)
	if err != nil {
		fmt.Printf("error during binding query, %s", err)
		return nil, &Error{Type: ErrorTypeClient, Message: "error during binding query", Details: err}
	}
	fmt.Printf("query, %d\n", query.Page)

	transactions, err := b.Controller.Service.Transaction.GetAll(query.Page)
	if err != nil {
		switch err.(type) {
		//case service.LessPointerZeroError:
		//	fmt.Printf("less pointer zero exception, %s\n", err)
		//	return
		default:
			fmt.Printf("error on get all transaction, %s\n", err)
			return nil, &Error{Type: ErrorTypeServer, Message: "error during getting all transactions", Details: err}
		}
	}

	fmt.Println("Successfully got all transactions")
	return getTransactionsResponse{Transactions: transactions}, nil
}

type getTransactionByFilterParams struct {
	Filter string `uri:"filter" json:"filter" binding:"required"`
}

type getTransactionByFilterResponse struct {
	Trans []entities.Transaction
}

func (b *blockChainController) getTransactionByFilter(c *gin.Context) (interface{}, *Error) {
	var body getTransactionByFilterParams

	err := c.ShouldBindUri(&body)
	if err != nil {
		fmt.Printf("error on get transaction by filter, %s\n", err)
		return nil, &Error{Type: ErrorTypeClient, Message: "error on get transactions by filter", Details: err}
	}
	fmt.Printf("body: %s\n", body)

	transactions, err := b.Service.Transaction.GetByFilter(body.Filter)
	if err != nil {
		switch err.(type) {
		case service.NilPointerDataError:
			fmt.Printf("nil pointer data error, %s\n", err)
			return nil, &Error{Type: ErrorTypeClient, Message: "wrong filter for transaction", Details: err}
		default:
			fmt.Printf("error during get data by filter, %s\n", err)
			return nil, &Error{Type: ErrorTypeServer, Message: "error during getting transaction", Details: err}
		}
	}

	fmt.Println("Successfully got transaction(s) by filter")
	return getTransactionByFilterResponse{Trans: transactions.Trans}, nil
}
