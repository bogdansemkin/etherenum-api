package app

import (
	"etherenum-api/etherenum-service/api/internal/config"
	httpController "etherenum-api/etherenum-service/api/internal/controller/http"
	"etherenum-api/etherenum-service/api/internal/repos"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/commission"
	"etherenum-api/etherenum-service/api/pkg/database"
	"etherenum-api/etherenum-service/api/pkg/etherscan"
	"etherenum-api/etherenum-service/api/pkg/hex"
	"etherenum-api/etherenum-service/api/pkg/logger"
	httpserver "etherenum-api/etherenum-service/api/pkg/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(config *config.Config, port string) error {
	logger := logger.NewZapLogger(config.Log.Level)

	collection, err := database.NewMongo(database.MongoDBConfig{
		Name:   config.Mongo.Name,
		User:   config.Mongo.User,
		Pass:   config.Mongo.Password,
		DBname: config.Mongo.DBname,
		//Host: config.Mongo.Host,
		//Port: config.Mongo.Port,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("error during creating mongoDB connection, %s", err))
	}
	converter := hex.NewConverter(logger)
	calculator := commission.NewCommissionCalculator(converter)
	repository := service.Repos{Transactions: repos.NewTransactionRepo(collection)}
	services := service.Service{Transaction: service.NewTransactionService(repository, logger, calculator)}
	etherscanner := etherscan.NewEtherscan(config, logger, services, converter)

	router := gin.New()

	httpController.New(httpController.Options{
		Handler: router,
		Config:  config,
		Service: services,
		Repos:   repository,
		Logger:  logger,
		Converter: converter,
	})
	err = etherscanner.InitBlocks()
	if err != nil {
		return fmt.Errorf("error during init blocks")
	}
	go func() {
		for {
			err := etherscanner.InputTransactions()
			if err != nil {
				fmt.Errorf("error during handling transactions, %s", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	httpServer := httpserver.New(
		router,
		httpserver.Port(port),
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
