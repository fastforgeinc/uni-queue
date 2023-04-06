package queue

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Options struct {
	Timeout time.Duration
}

type ListQueue struct {
	Redis rediser
	opt   *Options
}

func NewListQueue(redis rediser, opt *Options) *ListQueue {
	queue := &ListQueue{
		Redis: redis,
		opt:   opt,
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
	return q.Redis.Del(ctx, queue).Err()
}

func (q *ListQueue) Len(ctx context.Context, queue string) int {
	return int(q.Redis.LLen(ctx, queue).Val())
}

func (q *ListQueue) enqueue(ctx context.Context, queue string, value interface{}) error {
	b, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return q.Redis.LPush(ctx, queue, b).Err()
}

func (q *ListQueue) dequeue(ctx context.Context, queue string, value interface{}) error {
	cmd := q.Redis.BRPop(ctx, q.opt.Timeout, queue)
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
