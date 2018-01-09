package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var (
	port *int
	n    uint64
)

func main() {
	port = flag.Int("l", 1111, "local port to listen")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()
		serveUDP(ctx, *port)
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	cancel()

	wg.Wait()

	log.Println("quit")

	// runtime.GC()
	// debug.FreeOSMemory()
	// f, _ := os.Create("heap.prof")
	// pprof.WriteHeapProfile(f)
}

func serveUDP(ctx context.Context, port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("listenning on %s\n", addr)

	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		for {
			before := atomic.LoadUint64(&n)
			time.Sleep(time.Second)
			if n := atomic.LoadUint64(&n) - before; n != 0 {
				log.Printf("%d/s\n", n)
			}
		}
	}()

	go func() {
		for {
			buffer := make([]byte, 1)
			_, _, err := l.ReadFrom(buffer)
			if err != nil {
				log.Printf("error: %s\n", err)
				continue
			}
			atomic.AddUint64(&n, 1)
		}
	}()

	<-ctx.Done()
	return
}
