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
	var phpUpstream = common.MiddlewareReverseProxy("https://httpbin.org")

	app.Engine.GET("/healthz", common.HealthzHandler())
	app.Engine.GET("/api/v1/version", common.VersionHandler())
	app.Engine.Any("status/:code", common.SetCanaryUpstream, common.MiddlewareCanary(map[string]gin.HandlerFunc{
		"go": func(c *gin.Context) {
			c.Header("X-Upstream", "go")
			c.Status(200)
			c.Abort()
		},
		"php": phpUpstream,
	}))

	app.Engine.Any("/api/v1/echo", Echo(app))

	app.Engine.POST("/api/v1/books", CreateBook(app))
	app.Engine.GET("/api/v1/books", ListBooks(app))
}
