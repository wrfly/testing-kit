package main

import (
	"context"
	"log"
	"time"

	"github.com/wrfly/testing-kit/util/tokenbucket"
)

func main() {
	var take, drop = 0, 0
	var start time.Time

	bkt := tokenbucket.New(1000, time.Second)

	log.Println("5s test with 1000/s")
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
	log.Println("used:", time.Since(start).Seconds(), "s")
	log.Println("take:", take)
	log.Println("drop:", drop)

	num := int(1e8)
	log.Printf("range [%d] test\n", num)
	take, drop = 0, 0
	start = time.Now()
	for i := 0; i < num; i++ {
		if bkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("used:", time.Since(start).Seconds(), "s")
	log.Println("take:", take)
	log.Println("drop:", drop)

	smothBkt := tokenbucket.NewSmoth(1000, time.Second)

	log.Println("[smoth bucket] 5s test with 1000/s")
	ctx, cancel = context.WithDeadline(context.Background(),
		time.Now().Add(time.Second*5))
	defer cancel()
	start = time.Now()
	for ctx.Err() == nil {
		if smothBkt.TakeOne() {
			take++
			continue
		}
		drop++
	}
	log.Println("used:", time.Since(start).Seconds(), "s")
	log.Println("take:", take)
	log.Println("drop:", drop)

	num = int(1e8)
	log.Printf("[smoth bucket] range [%d] test\n", num)
	take, drop = 0, 0
	start = time.Now()
	for i := 0; i < num; i++ {
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
