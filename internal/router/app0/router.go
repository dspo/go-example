package app0

import (
	"net/http"

	"gitee.com/huajinet/go-example/internal/common"
)

func NewRouter(app ApplicationContext) http.Handler {
	mount(app)
	return app.Engine
}

func mount(app ApplicationContext) {
	app.Engine.GET("/healthz", common.HealthzHandler())
	app.Engine.GET("/api/v1/version", common.VersionHandler())

	app.Engine.Any("/api/v1/echo", Echo(app))
}
