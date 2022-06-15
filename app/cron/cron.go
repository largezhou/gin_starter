package cron

import (
	"github.com/largezhou/gin_starter/app"
	"github.com/largezhou/gin_starter/app/helper"
	"github.com/largezhou/gin_starter/app/logger"
)

func init() {
	app.NewCron("*/5 * * * * ?", func() {
		ctx := helper.NewTraceIdContext()
		logger.Debug(ctx, "cron")
	}).SkipIfStillRunning()
}
