package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type TestStruct struct {
	Bool    bool
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	Bytes   []byte
	String  string
}

func setupTest() (func(), context.Context, *redis.Client, *ListQueue) {
	const addr = "localhost:6379"
	ctx := context.Background()

	log.Println("setup suite")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	queue := NewListQueue(client)

	return func() {
		log.Println("teardown test")
	}, ctx, client, queue
}

func TestListQueue(t *testing.T) {
	const key = "test"

	teardownTest, ctx, client, queue := setupTest()
	defer teardownTest()

	t.Run("client_ping", func(t *testing.T) {
		cmd := client.Ping(ctx)
		assert.NoError(t, cmd.Err())
	})

	t.Run("queue_initial_flush", func(t *testing.T) {
		err := queue.Flush(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("len_before_enqueue", func(t *testing.T) {
		l := queue.Len(ctx, key)
		assert.Exactly(t, 0, l)
	})

	table := []struct {
		name  string
		value interface{}
	}{
		{"int32", int32(32)},
		{"int64", int64(64)},
		{"float32", float32(32.32)},
		{"float64", float64(64.64)},
		{"bytes", []byte{0xDE, 0xAD, 0xBE, 0xAF}},
		{"string", "foo"},
	}

	ts := TestStruct{true, 32, 64, 32.32, 64.64, []byte{0xC0, 0xCA, 0xC0, 0x1A}, "bar"}
	tsp := &ts

	for _, tc := range table {
		t.Run("enqueue_"+tc.name, func(t *testing.T) {
			err := queue.Enqueue(ctx, key, &tc.value)
			assert.NoError(t, err)
		})
	}

	t.Run("enqueue_struct", func(t *testing.T) {
		err := queue.Enqueue(ctx, key, &ts)
		assert.NoError(t, err)
	})

	t.Run("enqueue_struct_pointer", func(t *testing.T) {
		err := queue.Enqueue(ctx, key, &tsp)
		assert.NoError(t, err)
	})

	t.Run("len_after_enqueue", func(t *testing.T) {
		l := queue.Len(ctx, key)
		assert.Exactly(t, len(table)+2, l)
	})

	for _, tc := range table {
		t.Run("dequeue_"+tc.name, func(t *testing.T) {
			var v interface{}
			err := queue.Dequeue(ctx, key, &v)
			if assert.NoError(t, err) {
				assert.Exactly(t, tc.value, v)
			}
		})
	}

	t.Run("dequeue_struct", func(t *testing.T) {
		var v TestStruct
		err := queue.Dequeue(ctx, key, &v)
		if assert.NoError(t, err) {
			assert.Exactly(t, ts, v)
		}
	})

	t.Run("dequeue_struct", func(t *testing.T) {
		var v *TestStruct
		err := queue.Dequeue(ctx, key, &v)
		if assert.NoError(t, err) {
			assert.Exactly(t, tsp, v)
		}
	})

	t.Run("len_after_dequeue", func(t *testing.T) {
		l := queue.Len(ctx, key)
		assert.Exactly(t, 0, l)
	})

	t.Run("queue_finish_flush", func(t *testing.T) {
		err := queue.Flush(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("client_close", func(t *testing.T) {
		err := client.Close()
		assert.NoError(t, err)
	})
}
