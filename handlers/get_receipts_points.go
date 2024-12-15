package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Receipt_id struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func GetReceiptsPoints(c *gin.Context) {
	var receipt_id Receipt_id
	if err := c.ShouldBindUri(&receipt_id); err != nil {
		log.Printf("Error binding URL: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "error", "message": "Invalid receipt id"})
		return
	}

	id, err := url.QueryUnescape(c.Param("id"))
	if err != nil {
		log.Printf("Error unescaping URL: %v", err)
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
