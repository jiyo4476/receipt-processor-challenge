package router

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jiyo4476/receipt-processor-challenge/handlers"
	"github.com/jiyo4476/receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SetUpRouter() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom format
		return fmt.Sprintf("[receipt-processor-server] %s - [%s] \"%s%s\033[0m %s %s %s%d\033[0m %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.MethodColor(),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCodeColor(),
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

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
