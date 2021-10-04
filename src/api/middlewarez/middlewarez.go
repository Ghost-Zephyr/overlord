/*
Package middlewarez for the overlord software.
This package contains middleware functions
for the overlord HTTP REST API.
*/
package middlewarez

import (
	"fmt"

	"bitsnthings.dev/overlord/src/state/conf"
	"github.com/savsgio/atreugo/v11"
)

func SetMiddlewares(server *atreugo.Atreugo) {
	server.UseAfter(ServerHeaders)
}

func ServerHeaders(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.Add("Server", fmt.Sprintf("OverLord/%s", conf.Version))
	return ctx.Next()
}
