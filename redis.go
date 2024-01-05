package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Rediser interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
	LLen(ctx context.Context, key string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}
