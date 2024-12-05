package main

import (
	"log"
	"net/http"
	"net/url"
	"receipt-processor-challenge/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Map for in-memory data storage
var memory_cache = make(map[string]models.Receipt)

// Returns an ID for the receipt
func processReceipt(c *gin.Context) {
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

func getReceiptsPoints(c *gin.Context) {
	id, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		log.Printf("Error unescaping URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "error", "message": "Invalid receipt id"})
		return
	}

	_, err = uuid.Parse(id)
	if err != nil {
		log.Printf("Error Parsing ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "error", "message": "Invalid receipt id"})
		return
	}

	receipt, ok := memory_cache[id]
	if ok {
		points, _ := receipt.Points()
		c.JSON(http.StatusOK, gin.H{
			"points": points,
		})
	} else {
		log.Printf("No receipt found for that id: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
	}
}
