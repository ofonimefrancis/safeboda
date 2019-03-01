package common

import "github.com/gin-gonic/gin"

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
