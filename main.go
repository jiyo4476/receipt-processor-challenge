package main

import (
	"net/http"

	"receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Map for in-memory data storage
var memory_cache = make(map[string]models.Receipt)

// Returns an ID for the receipt
func processReceipt(c *gin.Context) {
	var receipt models.Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validate.Struct(receipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id = uuid.New().String()
	memory_cache[id] = receipt
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getReceiptsPoints(c *gin.Context) {
	receipt, ok := memory_cache[c.Param("id")]
	if ok {
		points := receipt.Points()
		c.JSON(http.StatusOK, gin.H{
			"points": points,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"code": "ID_NOT_FOUND", "message": "No receipt found for that id"})
	}
}

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("correctRetailerName", models.CorrectRetailerName)
		v.RegisterValidation("correctShortDescription", models.CorrectShortDescription)
		v.RegisterValidation("correctCashValue", models.CorrectCashValue)
	}

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptsPoints)

	router.Run()
}
