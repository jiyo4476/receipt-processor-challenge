package middleware

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ReceiptPointsLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()
		zap.L().Info("Retrieving points",
			zap.String("request_id", requestid.Get(c)))

		c.Next()

		// after request
		latency := time.Since(t)
		status := c.Writer.Status()

		zap.L().Info("Finished Retrieving points",
			zap.String("request_id", requestid.Get(c)),
			zap.Int("status", status),
			zap.String("latency", latency.String()),
		)

	}
}

func ProcessReceiptLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()
		zap.L().Info("Processing Receipt",
			zap.String("request_id", requestid.Get(c)))

		c.Next()

		// after request
		latency := time.Since(t)
		status := c.Writer.Status()

		zap.L().Info("Processed Request",
			zap.String("request_id", requestid.Get(c)),
			zap.Int("status", status),
			zap.String("latency", latency.String()),
		)

	}
}
