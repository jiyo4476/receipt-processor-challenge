package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jiyo4476/receipt-processor-challenge/models"
	"github.com/jiyo4476/receipt-processor-challenge/spec"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
func SetUpRouter() *gin.Engine {
	router := gin.Default()

	// Register custom validation functions for the test router
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("correctRetailerName", models.CorrectRetailerName)
		v.RegisterValidation("correctShortDescription", models.CorrectShortDescription)
		v.RegisterValidation("correctCashValue", models.CorrectCashValue)
		v.RegisterValidation("correctDate", models.CorrectDate)
		v.RegisterValidation("correctTime", models.CorrectTime)
	}

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptsPoints)
	return router
}

func main() {
	// Load specs in globally accessible variable
	if err := spec.PrintSpec("api.yml"); err != nil {
		log.Printf("Error loading spec: %v", err) // Log the error
		return
	}

	router := SetUpRouter()

	value, exists := os.LookupEnv("PORT")
	if !exists {
		value = "8080"
	}

	server := &http.Server{
		Addr:    ":" + value,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("received interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}
