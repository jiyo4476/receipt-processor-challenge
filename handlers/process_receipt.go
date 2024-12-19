package handlers

import (
	"net/http"

	"github.com/jiyo4476/receipt-processor-challenge/models"
	"github.com/jiyo4476/receipt-processor-challenge/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Returns an ID for the receipt
func ProcessReceipt(c *gin.Context) {
	zap.L().Info("Processing receipt")
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {

		zap.L().Warn("Validation Error",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to process receipt",
			"message": err.Error(),
		})
		return
	}

	var id = uuid.New().String()
	store.Receipts[id] = receipt
	zap.L().Info("Added receipt %s to database",
		zap.String("receipt_id", id),
	)
	c.JSON(http.StatusOK, gin.H{"id": id})
}
