package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/lib/config"
	"github.com/largezhou/gin_starter/lib/helper"
	"github.com/largezhou/gin_starter/lib/logger"
	"go.uber.org/zap"
)

var c = config.Config.App
var r *gin.Engine

func main() {
	defer helper.CallShutdownFunc()

	r = InitServer()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("造成恐慌")
	})

	srv := &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: r,
	}

	go func() {
		logger.Info("开始运行", zap.String("host", c.Host), zap.String("port", c.Port))
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

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
