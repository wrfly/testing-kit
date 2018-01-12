package main

import (
	"context"
	"log"
	"time"

	"github.com/wrfly/testing-kit/utils/tokenbucket"
)

func main() {
	var take, drop = 0, 0
	var start time.Time

	bkt := tokenbucket.New(1000, time.Second)

	log.Println("5s test with 1000/s")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		start = time.Now()
		for ctx.Err() == nil {
			if bkt.TakeOne() {
				take++
				continue
			}
			drop++
		}
		log.Println("used:", time.Since(start).Nanoseconds(), "ns")
		log.Println("take:", take)
		log.Println("drop:", drop)
	}()
	time.Sleep(5 * time.Second)
	cancel()

	time.Sleep(time.Millisecond)

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
	log.Println("used:", time.Since(start).Nanoseconds(), "ns")
	log.Println("take:", take)
	log.Println("drop:", drop)
}
