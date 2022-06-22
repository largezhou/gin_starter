package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/largezhou/gin_starter/app/command"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/init_ctx"
	"github.com/largezhou/gin_starter/app/logger"
	"github.com/largezhou/gin_starter/app/middleware"
	cronPkg "github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var cfg = config.Config.App
var Engine *gin.Engine
var Console *cli.App
var Cron *cronPkg.Cron

// Args 去掉 命令行 cli 标识参数后的 运行参数
var Args []string

func init() {
	ctx := init_ctx.Ctx
	Engine = initServer(ctx)
	Console = &cli.App{
		Commands: command.Commands,
	}
	Args = initArgs(ctx)

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err)
	}

	cronLogger := CronLogger{}
	Cron = cronPkg.New(
		cronPkg.WithSeconds(),
		cronPkg.WithLocation(loc),
		cronPkg.WithChain(
			cronPkg.Recover(cronLogger),
			SkipIfStillRunning(cronLogger),
			DelayIfStillRunning(cronLogger),
			SkipIfRunningOnOtherServer(cronLogger),
		),
		cronPkg.WithLogger(cronLogger),
	)
}

func RunInConsole(ctx context.Context) bool {
	args := os.Args
	return len(args) >= 2 && args[1] == app_const.CliKey
}

func initServer(ctx context.Context) *gin.Engine {
	if RunInConsole(ctx) {
		return nil
	}

	c := config.Config.App
	if c.Debug {
		logger.Debug(ctx, "debug 模式运行")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.MaxMultipartMemory = 10 << 20
	r.Use(
		middleware.SetTraceId(),
		middleware.Cors(),
		middleware.Logger(),
	)

	return r
}

func initArgs(ctx context.Context) []string {
	var args []string
	// 第二个参数为 CLI 入口，复制并删除
	for i, arg := range os.Args {
		if i != 1 {
			args = append(args, arg)
		}
	}

	return args
}

func Run(ctx context.Context) {
	if RunInConsole(ctx) {
		runConsole(ctx)
	} else {
		runServer(ctx)
	}
}

func runServer(ctx context.Context) {
	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: Engine,
	}

	go func() {
		logger.Info(ctx, "开始运行", zap.String("host", cfg.Host), zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	if cfg.Cron {
		runCron(ctx)
		defer Cron.Stop()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "服务关闭出错", zap.Error(err))
	}

	logger.Info(ctx, "服务已关闭")
}

func runConsole(ctx context.Context) {
	if err := Console.Run(Args); err != nil {
		logger.Error(ctx, "命令行执行失败", zap.Error(err))
		panic(err)
	}
}

func runCron(ctx context.Context) {
	for _, cron := range cronList {
		if cron.runImmediate {
			cron.Run()
		}

		if _, err := Cron.AddJob(cron.spec, cron); err != nil {
			panic(err)
		}
	}
	Cron.Start()
}
