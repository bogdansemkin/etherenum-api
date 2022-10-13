package http

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config   *config.Config
}

type ControllerOptions struct {
	Handler  *gin.RouterGroup
	Config   *config.Config
}
type Options struct {
	Handler  *gin.Engine
	Config   *config.Config
}

func New(options Options) {
	_ = ControllerOptions{
		Handler:  options.Handler.Group("/api/v1"),
		Config:   options.Config,
	}
}
