package commands

import (
	"github.com/largezhou/gin_starter/app/init_ctx"
	"github.com/largezhou/gin_starter/app/model"
	"github.com/urfave/cli/v2"
)

var db = model.DB.WithContext(init_ctx.Ctx)

var Commands = []*cli.Command{
	NewMakeMigrationCommand(),
	NewMigrateInstallCommand(),
	NewMigrateRollbackCommand(),
}
