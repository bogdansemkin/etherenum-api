package app

import (
	"etherenum-api/etherenum-service/api/internal/config"
	httpController "etherenum-api/etherenum-service/api/internal/controller/http"
	"etherenum-api/etherenum-service/api/internal/repos"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/database"
	"etherenum-api/etherenum-service/api/pkg/etherscan"
	"etherenum-api/etherenum-service/api/pkg/logger"
	httpserver "etherenum-api/etherenum-service/api/pkg/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(config *config.Config) error {
	collection, err := database.NewMongo(database.MongoDBConfig{
		Name: config.Mongo.Name,
		Port: config.Mongo.Port,
		Host: config.Mongo.Host,
	})
	if err != nil {
		return fmt.Errorf("error during creating mongoDB connection, %s", err)
	}

	repository := service.Repos{Transactions: repos.NewTransactionRepo(collection)}
	services := service.Service{Transaction: service.NewTransactionService(repository)}
	etherscanner := etherscan.NewEtherscan(config)
	logger := logger.NewLogger()

	router := gin.New()

	httpController.New(httpController.Options{
		Handler: router,
		Config:  config,
		Service: services,
		Repos:   repository,
	})

	//TODO need to refactor
	go func() {
		for {
			func() {
				body, err := etherscanner.GetBlock()
				logs := logger.GetLogs()
				if body == nil {
					return
				}
				if logs[len(logs)-1] == body.Result {
					return
				}
				logger.CreateLog(body.Result)
				fmt.Println(logs)

				allTransactions, err := etherscanner.GetTransactions(body.Result)
				if err != nil {
					fmt.Printf("error during getting the transaction, %s\n", err)
					return
				}

				previousBlockTransactions, err := repository.Transactions.GetByFilter(logs[len(logs)-1])
				if err != nil {
					fmt.Printf("error during getting the transaction, %s\n", err)
				}
				var transactionz []etherscan.Transaction
				for i := range previousBlockTransactions.Trans {
					transactionz = append(transactionz, etherscan.Transaction{
						Blockhash:        previousBlockTransactions.Trans[i].Blockhash,
						BlockNumber:      previousBlockTransactions.Trans[i].BlockNumber,
						From:             previousBlockTransactions.Trans[i].From,
						Gas:              previousBlockTransactions.Trans[i].Gas,
						GasPrice:         previousBlockTransactions.Trans[i].GasPrice,
						Hash:             previousBlockTransactions.Trans[i].Hash,
						Input:            previousBlockTransactions.Trans[i].Input,
						Nonce:            previousBlockTransactions.Trans[i].Nonce,
						To:               previousBlockTransactions.Trans[i].To,
						TransactionIndex: previousBlockTransactions.Trans[i].TransactionIndex,
						ChainId:          previousBlockTransactions.Trans[i].ChainId,
						AcceptNumber:     previousBlockTransactions.Trans[i].AcceptNumber,
					})
				}
				factoredTransactions, err := etherscanner.AcceptIncrement(allTransactions, transactionz)

				if factoredTransactions == nil {
					factoredTransactions = allTransactions
				}
				var trainers []interface{}
				for i := range factoredTransactions {
					//fmt.Printf("Transactions hash: %s, iterations: %d\n",factoredTransactions[i].Hash,factoredTransactions[i].AcceptNumber)
					trainers = append(trainers, factoredTransactions[i])
				}
				err = repository.Transactions.Insert(trainers)
				if err != nil {
					fmt.Errorf("error during inserting, %s", err)
					return
				}
			}()
			time.Sleep(400 * time.Millisecond)
		}
	}()

	httpServer := httpserver.New(
		router,
		httpserver.Port(config.HTTP.Port),
		httpserver.ReadTimeout(time.Second*60),
		httpserver.WriteTimeout(time.Second*60),
		httpserver.ShutdownTimeout(time.Second*30),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Printf("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		fmt.Printf("app - Run - httpServer.Notify. %s", err)
	}

	// shutdown http server
	err = httpServer.Shutdown()
	if err != nil {
		fmt.Printf("app - Run - httpServer.Shutdown. %s", err)
	}

	return nil
}
