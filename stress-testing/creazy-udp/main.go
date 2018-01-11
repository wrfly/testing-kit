package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	remote    *string
	data      *string
	rate      *int
	port      *int
	clientNum *int

	connGroup []net.Conn

	sentPkt uint64 // send per second
)

func main() {
	cpuNum := runtime.NumCPU()

	rate = flag.Int("n", 10000, "send rate, <num>/s, set to -1 for ultimate")
	remote = flag.String("r", "localhost:12345", "remote addr to send packet")
	data = flag.String("d", "", "data to send, empty by default")
	clientNum = flag.Int("c", cpuNum, "client num, default is `runtime.NumCPU()`")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("send '%s' to [%s] with [%d] clients at %d/s", *data,
		*remote, *clientNum, *rate)

	for i := 0; i < *clientNum; i++ {
		conn, err := net.Dial("udp", *remote)
		if err != nil {
			log.Printf("conn error: %s", err)
			return
		}
		connGroup = append(connGroup, conn)
	}

	if *rate > 0 {
		go func() {
			numPerClient := *rate / *clientNum
			for {
				for _, conn := range connGroup {
					go func(conn net.Conn) {
						for i := 0; i < numPerClient; i++ {
							atomic.AddUint64(&sentPkt, 1)
							conn.Write([]byte(*data))
						}
					}(conn)
				}
				time.Sleep(time.Second)
			}
		}()
	} else {
		go func() {
			for {
				for _, conn := range connGroup {
					atomic.AddUint64(&sentPkt, 1)
					conn.Write([]byte(*data))
				}
			}
		}()
	}

	go statistic(ctx)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	cancel()

	for _, conn := range connGroup {
		conn.Close()
	}

	log.Println("quit")
}

func statistic(ctx context.Context) {
	s := atomic.LoadUint64(&sentPkt)
	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			s = atomic.LoadUint64(&sentPkt)
			atomic.StoreUint64(&sentPkt, 0)
			if s != 0 {
				log.Printf("send: %v/s\n", s)
			}
		}
	}
}
