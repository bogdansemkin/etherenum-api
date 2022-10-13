package http

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/internal/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config   *config.Config
	Service  service.Service
	Repos    service.Repos
}

type ControllerOptions struct {
	Handler  *gin.RouterGroup
	Config   *config.Config
	Service  service.Service
	Repos    service.Repos
}
type Options struct {
	Handler  *gin.Engine
	Config   *config.Config
	Service  service.Service
	Repos    service.Repos
}

func New(options Options) {
	routerOptions := ControllerOptions{
		Handler:  options.Handler.Group("/api/v1"),
		Service: options.Service,
		Repos:    options.Repos,
		Config:   options.Config,
	}
	NewBlockchainRoutes(routerOptions)
}
