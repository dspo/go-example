package app0

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitee.com/huajinet/go-example/internal/common"
)

func NewRouter(app ApplicationContext) http.Handler {
	mount(app)
	return app.Engine
}

func mount(app ApplicationContext) {
	app.Engine.GET("/healthz", common.HealthzHandler())
	app.Engine.GET("/v1/version", common.VersionHandler())

	app.Engine.Any("/v1/echo", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	})
}
