package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jiyo4476/receipt-processor-challenge/router"
	"github.com/jiyo4476/receipt-processor-challenge/spec"
	"github.com/kelseyhightower/envconfig"
)

var logger *zap.Logger
var resetLogger func()
var isProd bool

type environment struct {
	PORT     string `default:"8080"`
	HOSTNAME string `default:"localhost"`
}

type environment_gin struct {
	MODE string `default:"debug"`
}

func getEnv() environment {
	var env environment
	err := envconfig.Process("RECEIPT_PROCESSOR", &env)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return env
}

func getIsProd() bool {
	var env environment_gin
	err := envconfig.Process("GIN", &env)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	if env.MODE == "release" {
		return true
	}
	return false
}

func init() {
	logger = getLogger()
	isProd = getIsProd()
}

func getLogger() *zap.Logger {
	var zapLogger *zap.Logger
	if isProd {
		zapLogger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer zapLogger.Sync()
	resetLogger = zap.ReplaceGlobals(zapLogger)
	return zapLogger
}

func getServer() *http.Server {
	env := getEnv()

	cur_router := router.SetUpRouter()

	// Add middleware
	cur_router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := requestid.Get(c); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			if timestamp := c.Value("timestamp"); timestamp != nil {
				fields = append(fields, zap.String("timestamp", timestamp.(string)))
			}

			return fields
		}),
	}))
	cur_router.Use(ginzap.RecoveryWithZap(logger, true))

	server := &http.Server{
		Addr:    env.HOSTNAME + ":" + env.PORT,
		Handler: cur_router,
	}

	return server
}

func main() {
	// Reset logger at the end of the main function
	defer resetLogger()

	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		logger.Fatal("Error loading spec",
			zap.Error(err))
		return
	}

	server := getServer()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	go func() {
		<-quit
		logger.Info("received interrupt signal")
		if err := server.Close(); err != nil {
			logger.Warn(fmt.Sprintf("Error closing server: %s", err.Error()))
		}
	}()

	logger.Info("Server Listening",
		zap.String("address", server.Addr),
	)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("Server closed under request")
		} else {
			logger.Error("Server closed unexpectedly",
				zap.Error(err))
		}
	}

	logger.Info("Server exiting")
}
