package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/app_const"
	"github.com/largezhou/gin_starter/app/helper"
	"github.com/largezhou/gin_starter/app/logger"
	"github.com/largezhou/gin_starter/app/redis"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type CronLogger struct {
}

func (c CronLogger) Info(msg string, keysAndValues ...any) {
	logger.WithChannel(app_const.LogCron).Info(context.Background(), msg, zap.Any("keysAndValues", keysAndValues))
}

func (c CronLogger) Error(err error, msg string, keysAndValues ...any) {
	logger.WithChannel(app_const.LogCron).Error(context.Background(), msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

type CronJob struct {
	spec string
	f    func()

	skipIfStillRunning  bool   // 前面的任务还未执行结束时，是否跳过
	delayIfStillRunning bool   // 前面的任务还未执行结束时，是否延迟执行
	runImmediate        bool   // 启动时是否直接执行一次
	name                string // 启用 runOnOneServer 时，必须指定唯一名字
	runOnOneServer      bool   // 通过 redis 分布式锁，确保只会在一台机器上执行
}

var cronList []*CronJob

func NewCron(spec string, f func()) *CronJob {
	c := &CronJob{
		spec: spec,
		f:    f,

		skipIfStillRunning:  false,
		delayIfStillRunning: false,
	}
	cronList = append(cronList, c)

	return c
}

func (c *CronJob) Run() {
	c.f()
}

// SkipIfStillRunning 如果上一次任务还在执行，则跳过当前的任务
func (c *CronJob) SkipIfStillRunning() *CronJob {
	c.skipIfStillRunning = true
	return c
}

// DelayIfStillRunning 如果上一次任务还在执行，则延迟执行
func (c *CronJob) DelayIfStillRunning() *CronJob {
	c.delayIfStillRunning = true
	return c
}

// RunImmediate 启动时会执行一次
func (c *CronJob) RunImmediate() *CronJob {
	c.runImmediate = true
	return c
}

func (c *CronJob) SetName(name string) *CronJob {
	c.name = name
	return c
}

// RunOnOneServer 通过 redis 分布式锁，确保只会在一台机器上执行
func (c *CronJob) RunOnOneServer() *CronJob {
	if c.name == "" {
		panic(errors.New("启用 RunOnOneServer 时，必须调用 SetName 设置一个名字"))
	}

	c.runOnOneServer = true
	return c
}

func DelayIfStillRunning(logger cron.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		if cj, ok := j.(*CronJob); ok && cj.delayIfStillRunning {
			var mu sync.Mutex
			return cron.FuncJob(func() {
				start := time.Now()
				mu.Lock()
				defer mu.Unlock()
				if dur := time.Since(start); dur > time.Minute {
					logger.Info("delayIfStillRunning", "duration", dur)
				}
				j.Run()
			})
		} else {
			return j
		}
	}
}

func SkipIfStillRunning(logger cron.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		if cj, ok := j.(*CronJob); ok && cj.skipIfStillRunning {
			var ch = make(chan struct{}, 1)
			ch <- struct{}{}
			return cron.FuncJob(func() {
				select {
				case v := <-ch:
					defer func() {
						// 避免 定时任务偶现的 panic 后，后续的任务全部会跳过
						if err := recover(); err != nil {
							if len(ch) == 0 {
								ch <- v
							}

							panic(err)
						}
					}()

					j.Run()

					ch <- v
				default:
					logger.Info("skipIfStillRunning")
				}
			})
		} else {
			return j
		}
	}
}

func SkipIfRunningOnOtherServer(logger cron.Logger) cron.JobWrapper {
	owner := uuid.NewString()
	return func(j cron.Job) cron.Job {
		if cj, ok := j.(*CronJob); ok && cj.runOnOneServer {
			return cron.FuncJob(func() {
				ctx := helper.NewTraceIdContext()
				lock := redis.NewLock(cj.name, 24*time.Hour, owner)
				defer lock.Unlock(ctx)
				if lock.TryLock(ctx) {
					j.Run()
				} else {
					logger.Info("skipIfRunningOnOtherServer")
				}
			})
		} else {
			return j
		}
	}
}
