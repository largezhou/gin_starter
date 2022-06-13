package redis

import (
	r "github.com/go-redis/redis/v8"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/init_ctx"
)

var cfg = config.Config.Redis
var Client *r.Client

func init() {
	Client = r.NewClient(&r.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.Db,
	}).WithContext(init_ctx.Ctx)
}
