package helper

import (
	"context"

	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/appconst"
)

// NewTraceIdContext 返回一个新的带链路追踪 ID 的 context
func NewTraceIdContext() context.Context {
	return context.WithValue(context.Background(), appconst.TraceIdKey, uuid.NewString())
}

// P 返回一个值的指针
func P[T any](v T) *T {
	return &v
}

// V 返回一个指针的底层值，指针为 nil 则返回 0 值
func V[T any](p *T) T {
	if p == nil {
		var t T
		return t
	}

	return *p
}
