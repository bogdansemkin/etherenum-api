package app

import (
	"etherenum-api/etherenum-service/api/internal/config"
	httpController "etherenum-api/etherenum-service/api/internal/controller/http"
	"etherenum-api/etherenum-service/api/internal/repos"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/database"
	"etherenum-api/etherenum-service/api/pkg/etherscan"
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

	router := gin.New()

	httpController.New(httpController.Options{
		Handler: router,
		Config:  config,
		Service: services,
		Repos:   repository,
	})

	go func() {
		for {
			allTransactions, err := etherscanner.GetTransactions()
			if err != nil {
				fmt.Errorf("error during getting the transaction, %s", err)
				return
			}

			trainers := []interface{}{}
			for i := range allTransactions {
				trainers = append(trainers, allTransactions[i])
			}

			err = repository.Transactions.Insert(trainers)
			if err != nil {
				fmt.Errorf("error during inserting, %s", err)
				return
			}
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
