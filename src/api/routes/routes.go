package routes

import "github.com/savsgio/atreugo/v11"

func SetRoutes(srv *atreugo.Atreugo) {
	srv.GET("/", func(ctx *atreugo.RequestCtx) error {
		return ctx.JSONResponse([]string{
			"Welcome to the overlord API.",
			"We're building an awesome VM/container management system!",
		})
	})
}
