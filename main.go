package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/lib/config"
	"github.com/largezhou/gin_starter/lib/helper"
	"github.com/largezhou/gin_starter/lib/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var c = config.Config.App

func main() {
	defer helper.CallShutdownFunc()

	if c.Debug {
		logger.Debug("debug 模式运行")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: r,
	}

	go func() {
		logger.Info("开始运行", zap.String("host", c.Host), zap.String("port", c.Port))
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Error("运行出错", zap.Error(err))
		}
	}()

	if _, err := os.Create("./xxx/yyy/zzz.log"); err != nil {
		logger.Error("打开文件错误", zap.Error(err))
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务关闭出错", zap.Error(err))
	}

	logger.Info("服务已关闭")
}
