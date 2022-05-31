package main

import (
	_ "github.com/largezhou/gin_starter/lib/viper"
	"github.com/spf13/viper"
)

type config struct {
	App   app // 应用基础配置
	Mysql mysql
}

type app struct {
	Host  string // 监听 IP
	Port  string // 监听 端口
	Env   string // 环境
	Debug bool   // 是否开启 debug
}

type mysql struct {
	Host     string
	Port     string
	User     string
	Password string
}

var Config config

func init() {
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
