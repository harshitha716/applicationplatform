package pantheon

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	serverconfig "github.com/Zampfi/application-platform/services/api/config"
	"github.com/Zampfi/application-platform/services/api/helper"

	"github.com/gin-gonic/gin"
)

func RegisterPantheonRoutes(router *gin.RouterGroup, serverCtx *serverconfig.ServerConfig) {
	pantheonURL := serverCtx.Env.PantheonURL
	url, err := url.ParseRequestURI(pantheonURL)
	if err != nil {
		panic(fmt.Errorf("failed to parse pantheon URL: %w", err))
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	router.Any(RouteAIProxyPath, func(c *gin.Context) {
		// Update the request to reflect the target scheme and host
		c.Request.URL.Scheme = url.Scheme
		c.Request.URL.Host = url.Host
		c.Request.Header.Set(XForwardedHost, c.Request.Host)
		c.Request.Host = url.Host
		c.Request.URL.Path = c.Param(ProxyPathParam)

		helper.AddAuthHeaders(c.Request.Header, c)

		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
