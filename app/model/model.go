package model

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/largezhou/gin_starter/app/appconst"
	"github.com/largezhou/gin_starter/app/config"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var c = config.Config.Mysql
var appConfig = config.Config.App
var DB *gorm.DB

type Model struct {
	Id         uint      `gorm:"primaryKey" json:"id"`
	CreateTime time.Time `gorm:"type:datetime;autoCreateTime;not null" json:"createTime"`
	UpdateTime time.Time `gorm:"type:datetime;autoUpdateTime;not null" json:"updateTime"`
}

type SqlRecorderLogger struct{}

func (l *SqlRecorderLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *SqlRecorderLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// ignore
}

func (l *SqlRecorderLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// ignore
}

func (l *SqlRecorderLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// ignore
}

func (l *SqlRecorderLogger) Trace(
	ctx context.Context,
	_ time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, rows := fc()
	logger.WithChannel(appconst.LogSql).Info(
		ctx,
		"sql",
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Error(err),
	)
}

func init() {
	dsn := c.Dsn
	if !strings.Contains(dsn, "loc=") {
		dsn += "&loc=" + url.QueryEscape(appConfig.Timezone)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &SqlRecorderLogger{},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
}
