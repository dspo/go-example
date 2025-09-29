package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VersionHandler returns the version of the apps
// todo: hard code yet, please implement it
func VersionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "0.0.1",
		})
	}
}
