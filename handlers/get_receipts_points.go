package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetReceiptsPoints(c *gin.Context) {
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
