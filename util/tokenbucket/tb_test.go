package tokenbucket

import (
	"context"
	"log"
	"testing"
	"time"
)

type tb interface {
	TakeOne() bool
}

func testBucket(bkt tb) {
	var take, drop = 0, 0
	var start time.Time

	log.Println("5s test")
	ctx, cancel := context.WithDeadline(context.Background(),
		time.Now().Add(time.Second*5))
	defer cancel()

	start = time.Now()
	for ctx.Err() == nil {
		if bkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("used:", time.Since(start), "s")
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	time.Sleep(time.Second)

	take = 0
	drop = 0
	end := int(1e8)
	log.Println("range test: ", end)
	start = time.Now()
	for i := 0; i < end; i++ {
		take++
		drop++
	}
	log.Println("dry run: ", time.Since(start).Seconds())

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if bkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("used:", time.Since(start).Seconds(), "s")
	log.Println("take:", take)
	log.Println("drop:", drop)
}

func TestSmothBucket(t *testing.T) {
	bkt := NewSmoth(1000, time.Second)
	testBucket(bkt)
}

func TestBucket(t *testing.T) {
	bkt := New(1000, time.Second)
	testBucket(bkt)
}
