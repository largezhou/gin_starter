package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/largezhou/gin_starter/app/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 记录请求，添加 requestId
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// 该类型特殊处理，因为可能有文件
		var data zap.Field
		if ctx.ContentType() != binding.MIMEMultipartPOSTForm {
			rawData, _ := ctx.GetRawData()
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
			data = zap.ByteString("data", rawData)
		} else {
			form, _ := ctx.MultipartForm()
			data = zap.Any("data", form)
		}

		logger.Info(
			ctx,
			"request",
			zap.String("clientIp", ctx.ClientIP()),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("query", ctx.Request.URL.RawQuery),
			data,
		)

		blw := &bodyLogWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = blw

		ctx.Next()

		logger.Info(
			ctx,
			"response",
			zap.Duration("cost", time.Now().Sub(start)),
			zap.Int("code", ctx.Writer.Status()),
			zap.ByteString("data", blw.body.Bytes()),
		)
	}
}
