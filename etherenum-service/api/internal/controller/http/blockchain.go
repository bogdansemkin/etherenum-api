package http

import (
	"etherenum-api/etherenum-service/api/internal/entities"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		p.GET("/", r.getTransactions)
		p.GET("/:filter", r.getTransactionByFilter)
	}
}

type getTransactionsQuery struct {
	Page string `form:"page" json:"page"`
}

type getTransactionsResponse struct {
	Transactions *[]entities.Transaction
}

func (b *blockChainController) getTransactions(c *gin.Context) {
	var query getTransactionsQuery
	 err := c.BindQuery(&query)
	 if err != nil {
	 	fmt.Printf("error during binding query, %s", err)
		 return
	 }
	 fmt.Printf("query, %s\n", query.Page)

	transactions, err := b.Controller.Service.Transaction.GetAll(query.Page)
	if err != nil {
		fmt.Printf("error on get all transaction, %s\n", err)
		return
	}

	//изменить if на switch-case
	if transactions == nil {
		c.JSON(http.StatusOK, "Hello world")
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: transactions})
}

type getTransactionByFilterParams struct {
	Filter string `uri:"filter" json:"filter" binding:"required"`
}

func (b *blockChainController) getTransactionByFilter(c *gin.Context) {
	var body getTransactionByFilterParams

	err := c.ShouldBindUri(&body)
	if err != nil {
		fmt.Printf("error on get transaction by filter, %s\n", err)
		return
	}
	fmt.Printf("body: %s\n", body)
	transactions, err := b.Service.Transaction.GetByFilter(body.Filter)
	if err != nil {
		fmt.Printf("error during get data by filter, %s", err)
		return
	}
	//TODO добавить switch с кастомной обработкой ошибок, таких, как пустая data
	c.JSON(http.StatusOK, transactions.Trans)
}
