package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Receipt_id struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func GetReceiptsPoints(c *gin.Context) {
	var receipt_id Receipt_id
	if err := c.ShouldBindUri(&receipt_id); err != nil {
		log.Printf("Error binding URL: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": err.Error()})
		return
	}

	id := c.Param("id")
	receipt, ok := memory_cache[id]
	if !ok {
		log.Printf("No receipt found for id: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
	}

	points, _ := receipt.Points()
	c.JSON(http.StatusOK, gin.H{
		"points": points,
	})
}
