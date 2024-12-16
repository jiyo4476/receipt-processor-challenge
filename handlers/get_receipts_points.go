package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type receipt_id struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func GetReceiptsPoints(c *gin.Context) {
	var receiptId receipt_id
	if err := c.ShouldBindUri(&receiptId); err != nil {
		log.Infof("Error binding id: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
		return
	}

	id := c.Param("id")
	receipt, ok := memory_cache[id]
	if !ok {
		log.Infof("No receipt found for id: %s", id)
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
	}

	points, _ := receipt.Points()
	c.JSON(http.StatusOK, gin.H{
		"points": points,
	})
}
