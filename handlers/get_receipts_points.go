package handlers

import (
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
		zap.L().Warn("Validation Error:",
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
		return
	}

	id := c.Param("id")

	receipt, ok := store.Receipts[id]
	if !ok {
		zap.L().Warn("No receipt for id",
			zap.String("receipt_id", id))
		c.JSON(http.StatusNotFound, gin.H{"code": "error", "message": "No receipt found for that id"})
		return
	}

	points, _ := receipt.Points()
	zap.L().Info("Calculated points",
		zap.String("receipt_id", id),
		zap.Int("points", int(points)))

	c.JSON(http.StatusOK, gin.H{
		"points": points,
	})
}
