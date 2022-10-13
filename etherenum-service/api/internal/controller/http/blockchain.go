package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type blockChainController struct {
	Controller
}

func NewBlockchainRoutes(options ControllerOptions) {
	r := &blockChainController{Controller{
		Config: options.Config,
		Service: options.Service,
		Repos: options.Repos,
	}}

	p :=options.Handler.Group("/transactions")
	{
		p.GET("/", r.getTransactions)
	}
}

func (b *blockChainController) getTransactions(c *gin.Context) {
	fmt.Println("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	transactions, err := b.Controller.Service.Transaction.GetAll()
	if err != nil {
		fmt.Errorf("error on get all transaction, %s", err)
		return
	}

	//изменить if на switch-case
	if transactions == nil {
		c.JSON(http.StatusOK, "Hello world")
	}

	c.JSON(http.StatusOK, transactions)
}
