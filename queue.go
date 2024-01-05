package queue

import (
	"context"
)

type Queuer interface {
	Enqueue(ctx context.Context, queue string, value interface{}) error
	Dequeue(ctx context.Context, queue string, value interface{}) error
	Flush(ctx context.Context, queue string) error
	Len(ctx context.Context, queue string) int
}
