package common

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// MiddlewareReverseProxy creates a Gin middleware that acts as a reverse proxy.
func MiddlewareReverseProxy(targetURL string) gin.HandlerFunc {
	remote, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	return func(c *gin.Context) {
		// Modify the request before forwarding (optional)
		// For example, add/remove headers, modify the path, etc.
		c.Request.Header.Add("X-Forwarded-For", c.ClientIP())
		c.Request.Host = remote.Host // Important for some backend servers

		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort() // Abort further middleware execution after proxying
	}
}
