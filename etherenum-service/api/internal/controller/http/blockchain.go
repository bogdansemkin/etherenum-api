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
	}
}

type getTransactionsResponse struct {
	Transactions *[]entities.Transaction
}

func (b *blockChainController) getTransactions(c *gin.Context) {
	transactions, err := b.Controller.Service.Transaction.GetAll()
	if err != nil {
		fmt.Errorf("error on get all transaction, %s", err)
		return
	}

	//изменить if на switch-case
	if transactions == nil {
		c.JSON(http.StatusOK, "Hello world")
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: transactions})
}
