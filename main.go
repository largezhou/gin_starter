package main

import (
	"github.com/largezhou/gin_starter/app"
	"github.com/largezhou/gin_starter/app/api"
	"github.com/largezhou/gin_starter/app/init_ctx"
	"github.com/largezhou/gin_starter/app/shutdown"
)

func main() {
	ctx := init_ctx.Ctx

	defer func() {
		shutdown.CallShutdownFunc(ctx)
	}()

	if !app.RunInConsole(ctx) {
		api.InitRouter(app.Engine)
	}

	app.Run(ctx)
}
