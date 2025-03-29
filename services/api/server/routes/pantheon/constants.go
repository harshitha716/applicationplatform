package pantheon

import "fmt"

const XForwardedHost = "X-Forwarded-Host"
const ProxyPathParam = "proxyPath"

// Routes
var RouteAIProxyPath = fmt.Sprintf("/ai/*%s", ProxyPathParam)
