package pool

import (
	"github.com/fzzy/radix/redis"
	. "testing"
)

func TestPool(t *T) {
	pool, err := NewPool("tcp", "localhost:6379", 10)
	if err != nil {
		t.Fatal(err)
	}

	conns := make([]*redis.Client, 20)
	for i := range conns {
		if conns[i], err = pool.Get(); err != nil {
			t.Fatal(err)
		}
	}

	for i := range conns {
		pool.Put(conns[i])
	}

	pool.Empty()
}

func TestPoolPutWithNonEmptyPipeline(t *T) {
	pool, err := NewPool("tcp", "localhost:6379", 1)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := pool.Get()
	if err != nil {
		t.Fatal(err)
	}

	conn.Append("set", "mykey", "1")
	pool.Put(conn)

	newConn, newErr := pool.Get()
	if newErr != nil {
		t.Fatal(err)
	}

	if replyErr := newConn.GetReply().Err; replyErr != redis.PipelineQueueEmptyError {
		t.Fatal("Fresh connection from pool has non-empty pipeline")
	}

	pool.Empty()
}
