package common

import (
	"cmp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddlewareCanaryKeyFunc(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c)
		c.Next()
	}
}

func SetCanaryUpstream(c *gin.Context) {
	upstream := c.Request.Header.Get("X-Canary-Upstream")
	c.Set("X-Canary-Upstream", cmp.Or(upstream, "go"))
}

func MiddlewareCanary(upstreams map[string]gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if handlerFunc, ok := upstreams[c.MustGet("X-Canary-Upstream").(string)]; ok {
			handlerFunc(c)
			c.Abort()
			return
		}
		c.Status(http.StatusNotFound)
		c.Abort()
	}
}
