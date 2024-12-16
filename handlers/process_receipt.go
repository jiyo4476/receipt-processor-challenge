package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/jiyo4476/receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Returns an ID for the receipt
func ProcessReceipt(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		log.Infof("Error binding JSON: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": err.Error(),
		})
		return
	}
	// Validates the item total and the receipt total
	if err := receipt.ValidateTotal(); err != nil {
		log.Infof("Mismatching Receipt and Item Total: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Error",
			"message": "The receipt total does not match the sum of the item totals",
		})
		return
	}

	var id = uuid.New().String()
	memory_cache[id] = receipt
	c.JSON(http.StatusOK, gin.H{"id": id})
}
