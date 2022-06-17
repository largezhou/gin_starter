package redis

import (
	"context"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
	"time"
)

const releaseLockLuaScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`

type Lock struct {
	name  string
	owner string
	ttl   time.Duration
}

func NewLock(name string, ttl time.Duration, owner string) Lock {
	return Lock{
		name:  name,
		owner: owner,
		ttl:   ttl,
	}
}

func (l *Lock) TryLock(ctx context.Context) bool {
	res, err := Client.SetNX(ctx, l.name, l.owner, l.ttl).Result()
	if err != nil {
		logger.Error(ctx, "获取锁失败", zap.Error(err))
	}
	return res
}

func (l *Lock) Unlock(ctx context.Context) {
	Client.Eval(ctx, releaseLockLuaScript, []string{l.name}, l.owner)
}
