package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/jiyo4476/receipt-processor-challenge/handlers"
	"github.com/jiyo4476/receipt-processor-challenge/models"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	return logger
}

func SetUpRouter() *gin.Engine {
	//router := gin.Default()
	log := setupLogger()
	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	// Register custom validation functions for the test router
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("correctRetailerName", models.CorrectRetailerName)
		v.RegisterValidation("correctShortDescription", models.CorrectShortDescription)
		v.RegisterValidation("correctCashValue", models.CorrectCashValue)
		v.RegisterValidation("correctDate", models.CorrectDate)
		v.RegisterValidation("correctTime", models.CorrectTime)
	}

	router.POST("/receipts/process", handlers.ProcessReceipt)
	router.GET("/receipts/:id/points", handlers.GetReceiptsPoints)
	return router
}
