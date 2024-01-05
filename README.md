# Universal queue library for Golang

uni-queue library defines basic `Queue` interface that implements queues using:
- Redis [Lists](https://redis.io/docs/data-types/lists/) - `ListQueue`
- Redis [Streams](https://redis.io/docs/data-types/streams/) - `StreamQueue` (TBD)
- AWS [SQS](https://aws.amazon.com/sqs/) - `SQS` (TBD)

## Installation
```shell
go get github.com/ypopivniak/uni-queue
```

## Quickstart

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ypopivniak/uni-queue"
)

type Object struct {
	Str string
	Num int
}

func main() {
	ctx := context.TODO()

	// Construct ®Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Construct ListQueue
	q := queue.NewListQueue(client)

	// Construct ListQueue with dequeue timeout
	timeout := time.Second
	q := queue.NewListQueue(client, queue.WithDequeueTimeout(timeout))

	// Enqueue value
	input := Object{"foo", 69}
	err := q.Enqueue(ctx, "queue:name", &input)
	if err != nil {
		log.Fatal(err)
	}

	// Dequeue value
	var output Object
	err := queue.Dequeue(ctx, "queue:name", &output)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", output)
}
```
