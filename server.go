package main

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/lib/logger"
	"github.com/largezhou/gin_starter/lib/middleware"
)

func InitServer() *gin.Engine {
	if c.Debug {
		logger.Debug("debug 模式运行")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery())

	return r
}
