package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/jiyo4476/receipt-processor-challenge/router"
	"github.com/jiyo4476/receipt-processor-challenge/spec"
	"github.com/kelseyhightower/envconfig"
)

type environment struct {
	PORT     string `default:"8080"`
	HOSTNAME string `default:""`
}

func setUpLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	return logger
}

func main() {
	var env environment
	err := envconfig.Process("receipt_processor", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		log.Printf("Error loading spec: %v", err) // Log the error
		return
	}

	cur_router := router.SetUpRouter()
	cur_router.Use(ginlogrus.Logger(setUpLogger()), gin.Recovery())

	server := &http.Server{
		Addr:    env.HOSTNAME + ":" + env.PORT,
		Handler: cur_router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("received interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Closed:", err)
		}
	}()

	log.Infof("Listening on %s:%s", env.HOSTNAME, env.PORT)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpectedly")
		}
	}

	log.Println("Server exiting")
}
