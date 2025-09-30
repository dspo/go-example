package app0

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Echo(app ApplicationContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	}
}
