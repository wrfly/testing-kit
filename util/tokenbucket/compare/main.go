package main

import (
	"context"
	"log"
	"time"

	"github.com/bsm/ratelimit"
	jujuRatelimit "github.com/juju/ratelimit"
	"github.com/tsenart/tb"
	"github.com/wrfly/testing-kit/util/tokenbucket"
)

func main() {

	log.Println("start testing...")
	jujuBkt := jujuRatelimit.NewBucketWithRate(1000, 1000)
	bkt := ratelimit.New(1000, time.Second)
	wrflyBkt := tokenbucket.New(1000, time.Second)
	smothBkt := tokenbucket.New(1000, time.Second)
	tbBkt := tb.NewBucket(1000, time.Second)

	var (
		take, drop = 0, 0
		duration   = time.Second * 5
		start      time.Time
		ctx        context.Context
		cancel     context.CancelFunc
	)

	log.Println("5s test")

	drop = 0
	take = 0
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(duration))
	start = time.Now()
	for ctx.Err() == nil {
		if jujuBkt.TakeAvailable(1) == 1 {
			take++
			continue
		}
		drop++
	}
	log.Println("juju[lock]:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(duration))
	start = time.Now()
	for ctx.Err() == nil {
		if bkt.Limit() {
			drop++
			continue
		}
		take++
	}
	log.Println("bsm[atomic]:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(duration))
	start = time.Now()
	for ctx.Err() == nil {
		if wrflyBkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("wrfly:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(duration))
	start = time.Now()
	for ctx.Err() == nil {
		if smothBkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("smoth:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(duration))
	start = time.Now()
	for ctx.Err() == nil {
		if tbBkt.Take(1) == 1 {
			take++
			continue
		}
		drop++
	}
	cancel()
	log.Println("tb: ", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	// range
	end := int(1e7)
	log.Println("range test: ", end)
	start = time.Now()
	for i := 0; i < end; i++ {
		take++
		drop++
	}
	log.Println("dry run: ", time.Since(start).Seconds())
	log.Println()

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if jujuBkt.TakeAvailable(1) == 1 {
			take++
			continue
		}
		drop++
	}
	log.Println("juju[lock]:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if bkt.Limit() {
			drop++
			continue
		}
		take++
	}
	log.Println("bsm[atomic]:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if wrflyBkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("wrfly:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if smothBkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("smoth:", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()

	drop = 0
	take = 0
	start = time.Now()
	for i := 0; i < end; i++ {
		if tbBkt.Take(1) == 1 {
			take++
			continue
		}
		drop++
	}
	log.Println("tb: ", time.Since(start).Seconds())
	log.Println("take:", take)
	log.Println("drop:", drop)
	log.Println()
}
