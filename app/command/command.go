package command

import (
	"github.com/largezhou/gin_starter/app/initctx"
	"github.com/largezhou/gin_starter/app/model"
	"github.com/urfave/cli/v2"
)

var db = model.DB.WithContext(initctx.Ctx)

var Commands = []*cli.Command{
	NewMakeMigrationCommand(),
	NewMigrateInstallCommand(),
	NewMigrateRollbackCommand(),
}
