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

	"github.com/jiyo4476/receipt-processor-challenge/middleware"
	"github.com/jiyo4476/receipt-processor-challenge/router"
	"github.com/jiyo4476/receipt-processor-challenge/spec"
	"github.com/kelseyhightower/envconfig"
)

type environment struct {
	PORT     string `default:"8080"`
	HOSTNAME string `default:"localhost"`
}

func getEnv() environment {
	var env environment
	err := envconfig.Process("RECEIPT_PROCESSOR", &env)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	return env
}

func getLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error when creating logger %s", err.Error()))
	}
	defer logger.Sync() // flushes buffer, if any
	return logger
}

func getServer() *http.Server {
	env := getEnv()

	cur_router := router.SetUpRouter()

	// Add middleware
	cur_router.Use(requestid.New())
	logger := zap.L()

	cur_router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := requestid.Get(c); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			return fields
		}),
	}))
	cur_router.Use(ginzap.RecoveryWithZap(logger, true))

	cur_router.Use(middleware.RateLimiter)

	server := &http.Server{
		Addr:    env.HOSTNAME + ":" + env.PORT,
		Handler: cur_router,
	}

	return server
}

func main() {
	logger := getLogger()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		logger.Sugar().Fatalf("Error loading spec: %v", err)
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
			logger.Sugar().Warnf("Error closing server: %s", err.Error())
		}
	}()

	logger.Sugar().Info(fmt.Sprintf("Listening on %s", server.Addr))

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("Server closed under request")
		} else {
			logger.Sugar().Errorf("Server closed unexpectedly: %v", err)
		}
	}

	logger.Info("Server exiting")
}
