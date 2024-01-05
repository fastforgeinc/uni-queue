package queue

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Option func(*options)

type options struct {
	dequeueTimeoutSeconds time.Duration
}

// WithDequeueTimeout sets timeout of enqueue operation, default indefinitely
func WithDequeueTimeout(t time.Duration) Option {
	return func(o *options) {
		o.dequeueTimeoutSeconds = t
	}
}

type ListQueue struct {
	redis Rediser
	opt   *options
}

func NewListQueue(redis Rediser, opts ...Option) *ListQueue {
	queue := &ListQueue{
		redis: redis,
		opt:   new(options),
	}
	for _, opt := range opts {
		opt(queue.opt)
	}
	return queue
}

func (q *ListQueue) Enqueue(ctx context.Context, queue string, value interface{}) error {
	return q.enqueue(ctx, queue, value)
}

func (q *ListQueue) Dequeue(ctx context.Context, queue string, value interface{}) error {
	return q.dequeue(ctx, queue, value)
}

func (q *ListQueue) Flush(ctx context.Context, queue string) error {
	return q.redis.Del(ctx, queue).Err()
}

func (q *ListQueue) Len(ctx context.Context, queue string) int {
	return int(q.redis.LLen(ctx, queue).Val())
}

func (q *ListQueue) enqueue(ctx context.Context, queue string, value interface{}) error {
	b, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return q.redis.LPush(ctx, queue, b).Err()
}

func (q *ListQueue) dequeue(ctx context.Context, queue string, value interface{}) error {
	cmd := q.redis.BRPop(ctx, q.opt.dequeueTimeoutSeconds, queue)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	vs := cmd.Val()
	if len(vs) != 2 {
		value = nil
		return nil
	}
	v := []byte(vs[1])
	return msgpack.Unmarshal(v, value)
}
