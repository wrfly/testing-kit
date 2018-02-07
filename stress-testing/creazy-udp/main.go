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
	target    *string
	data      *string
	rate      *int
	port      *int
	clientNum *int
	ttnum     *uint64

	connGroup []net.Conn

	sentPkt uint64 // send per second
)

func main() {
	cpuNum := runtime.NumCPU()

	rate = flag.Int("r", 10000, "send rate, <num>/s, set to -1 for ultimate")
	ttnum = flag.Uint64("n", 0, "total number, set to positive for counter mode (ignore rate)")
	target = flag.String("t", "localhost:12345", "target addr to send packet")
	data = flag.String("d", "", "data to send, empty by default")
	clientNum = flag.Int("c", cpuNum, "client num, default is `runtime.NumCPU()`")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < *clientNum; i++ {
		conn, err := net.Dial("udp", *target)
		if err != nil {
			log.Printf("conn error: %s", err)
			return
		}
		connGroup = append(connGroup, conn)
	}

	if *rate > 0 && *ttnum == 0 {
		log.Printf("send '%s' to [%s] with [%d] clients at %d/s\n", *data,
			*target, *clientNum, *rate)
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
	} else if *rate < 0 && *ttnum == 0 {
		log.Printf("send '%s' to [%s] with [%d] clients in ultimate\n", *data,
			*target, *clientNum)
		go func() {
			for {
				for _, conn := range connGroup {
					atomic.AddUint64(&sentPkt, 1)
					conn.Write([]byte(*data))
				}
			}
		}()
	} else {
		log.Printf("send '%s' to [%s] with [%d] clients fot [%d] times\n", *data,
			*target, *clientNum, *ttnum)
		go func() {
			defer log.Printf("Send %d packets, done\n", *ttnum)
			for {
				for _, conn := range connGroup {
					atomic.AddUint64(&sentPkt, 1)
					conn.Write([]byte(*data))
					if atomic.LoadUint64(&sentPkt) >= *ttnum {
						return
					}
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
	log.Printf("total sent %v\n", atomic.LoadUint64(&sentPkt))
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
			if n := atomic.LoadUint64(&sentPkt) - s; n != 0 {
				s = atomic.LoadUint64(&sentPkt)
				log.Printf("Send: %v/s\tTotal sent %v\n", n, s)
			}
		}
	}
}
