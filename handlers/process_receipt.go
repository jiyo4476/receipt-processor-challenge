package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/jiyo4476/receipt-processor-challenge/models"
	"github.com/jiyo4476/receipt-processor-challenge/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Returns an ID for the receipt
func ProcessReceipt(c *gin.Context) {
	zap.L().Info(fmt.Sprintf("Processing receipt from request %s", requestid.Get(c)))
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		zap.L().Warn(fmt.Sprintf("Validation Error: %v", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to process receipt",
			"message": err.Error(),
		})
		return
	}

	var id = uuid.New().String()
	store.Receipts[id] = receipt
	zap.L().Info(fmt.Sprintf("Added receipt %s to database", id))
	c.JSON(http.StatusOK, gin.H{"id": id})
}
