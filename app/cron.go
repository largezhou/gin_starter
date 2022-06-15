package app

import (
	"context"
	"github.com/largezhou/gin_starter/app/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type CronLogger struct {
}

func (c CronLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Info(context.Background(), msg, zap.Any("keysAndValues", keysAndValues))
}

func (c CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Error(context.Background(), msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}

type CronJob struct {
	spec string
	f    func()

	skipIfStillRunning  bool
	delayIfStillRunning bool
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

func (c *CronJob) SkipIfStillRunning() *CronJob {
	c.skipIfStillRunning = true
	return c
}

func (c *CronJob) DelayIfStillRunning() *CronJob {
	c.delayIfStillRunning = true
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
