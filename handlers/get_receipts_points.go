package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jiyo4476/receipt-processor-challenge/store"
	"go.uber.org/zap"
)

type receipt_id struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func GetReceiptsPoints(c *gin.Context) {
	var receiptId receipt_id
	if err := c.ShouldBindUri(&receiptId); err != nil {
		zap.L().Warn(fmt.Sprintf("Validation Error:  %v", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
		return
	}

	id := c.Param("id")
	zap.L().Info(fmt.Sprintf("Getting points for %s", id))

	receipt, ok := store.Receipts[id]
	if !ok {
		zap.L().Warn(fmt.Sprintf("No receipt found for id: %s", id))
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
		return
	}

	points, _ := receipt.Points()
	zap.L().Info(fmt.Sprintf("%d points found for id %s", points, id))
	c.JSON(http.StatusOK, gin.H{
		"points": points,
	})
}
