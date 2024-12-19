package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 5)

// Middleware to check the rate limit.
func RateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		zap.L().Warn("To many requests")
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests please try again later"})
		c.Abort()
		return
	}
	c.Next()
}
