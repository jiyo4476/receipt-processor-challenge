package handlers

import (
	"log"
	"net/http"

	"github.com/jiyo4476/receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Returns an ID for the receipt
func ProcessReceipt(c *gin.Context) {
	var receipt models.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		log.Printf("Error binding JSON: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validates the item total and the receipt total
	if err := receipt.Validate(); err != nil {
		log.Printf("Validation error: %v", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id = uuid.New().String()
	memory_cache[id] = receipt
	c.JSON(http.StatusOK, gin.H{"id": id})
}
