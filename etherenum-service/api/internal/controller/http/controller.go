package http

import (
	"bytes"
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

type Controller struct {
	Config  *config.Config
	Service service.Service
	Repos   service.Repos
	Logger  logger.Logger
}

type ControllerOptions struct {
	Handler *gin.RouterGroup
	Config  *config.Config
	Service service.Service
	Repos   service.Repos
	Logger  logger.Logger
}
type Options struct {
	Handler *gin.Engine
	Config  *config.Config
	Service service.Service
	Repos   service.Repos
	Logger  logger.Logger
}

func New(options Options) {
	routerOptions := ControllerOptions{
		Handler: options.Handler.Group("/api/v1"),
		Service: options.Service,
		Repos:   options.Repos,
		Config:  options.Config,
		Logger:  options.Logger.Named("HTTPController"),
	}
	NewBlockchainRoutes(routerOptions)
}

type Error struct {
	Type          ErrorType   `json:"-"`
	Message       string      `json:"message"`
	Details       interface{} `json:"details,omitempty"`
	InvalidFields interface{} `json:"invalidFields,omitempty"`
}

// ErrorType is used to define error type.
type ErrorType string

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer ErrorType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient ErrorType = "client"
)

// Error is used to convert an error to a string.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func errorHandler(options ControllerOptions, handler func(c *gin.Context) (interface{}, *Error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// handle panics
		defer func() {
			if err := recover(); err != nil {
				// get stacktrace
				stacktrace, errors := gostackparse.Parse(bytes.NewReader(debug.Stack()))
				if len(errors) > 0 || len(stacktrace) == 0 {
					fmt.Errorf("get stacktrace errors, stacktraceErrors, errors: %s\n, stacktrace: %v ", err, stacktrace)
				} else {
					fmt.Errorf("unhandled error: %s, stacktrace: %v", err, stacktrace)
				}

				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					fmt.Errorf("failed to abort with error: %s", err)
				}
			}
		}()

		// execute handler
		body, err := handler(c)

		// check if middleware
		if body == nil && err == nil {
			return
		}
		fmt.Printf("Body %s, error: %s", body, err)

		if err != nil {
			if err.Type == ErrorTypeServer {
				fmt.Println("internal server error")

				// whether to send error to the client
				if options.Config.HTTP.SendDetailsOnInternalError {
					c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				} else {
					err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
					if err != nil {
						fmt.Errorf("failed to abort with error: %s", err)
					}
				}
			} else {
				fmt.Println("client error")
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			}
			return
		}

		fmt.Println("request handled")
		c.JSON(http.StatusOK, body)
	}
}
