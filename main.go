package main

import (
	"github.com/largezhou/gin_starter/app"
	"github.com/largezhou/gin_starter/app/api"
	_ "github.com/largezhou/gin_starter/app/cron"
	"github.com/largezhou/gin_starter/app/initctx"
	"github.com/largezhou/gin_starter/app/shutdown"
)

func main() {
	ctx := initctx.Ctx

	defer func() {
		shutdown.CallShutdownFunc(ctx)
	}()

	if !app.RunInConsole(ctx) {
		api.InitRouter(app.Engine)
	}

	app.Run(ctx)
}
