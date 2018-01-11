package main

import (
	"context"
	"flag"
	"fmt"
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
	times     *int
	port      *int
	clientNum *int

	connGroup []net.Conn

	sentPkt    uint64 // send per second
	receivePkt uint64 // receive per second
)

type packet struct {
	data []byte
}

func main() {
	cpuNum := runtime.NumCPU()

	port = flag.Int("l", 1234, "local port to listen")
	times = flag.Int("t", 10000, "how many times to repeat")
	remote = flag.String("r", "localhost:12345", "remote addr to send packet")
	clientNum = flag.Int("c", cpuNum, "client num, default is `runtime.NumCPU()`")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < *clientNum; i++ {
		conn, err := net.Dial("udp", *remote)
		if err != nil {
			log.Printf("conn error: %s", err)
			return
		}
		connGroup = append(connGroup, conn)
	}

	go func() {
		numPerClient := *times / *clientNum
		pktChan := serveUDP(ctx, *port)
		for pkt := range pktChan {
			go func(pkt packet) {
				for _, conn := range connGroup {
					go func(conn net.Conn) {
						for i := 0; i < numPerClient; i++ {
							atomic.AddUint64(&sentPkt, 1)
							conn.Write(pkt.data)
						}
					}(conn)
				}
			}(pkt)
		}
	}()

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
	r := atomic.LoadUint64(&receivePkt)
	s := atomic.LoadUint64(&sentPkt)
	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			r = atomic.LoadUint64(&receivePkt)
			s = atomic.LoadUint64(&sentPkt)
			atomic.StoreUint64(&receivePkt, 0)
			atomic.StoreUint64(&sentPkt, 0)
			if r != 0 || s != 0 {
				log.Printf("receive: %v/s\tsend: %v/s\n", r, s)
			}
		}
	}
}

func serveUDP(ctx context.Context, port int) chan packet {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("listenning on %s\n", addr)

	byteChan := make(chan packet, 1000)

	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 20480)
	go func() {
		for {
			size, _, err := l.ReadFrom(buffer)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Println(err)
				continue
			}
			atomic.AddUint64(&receivePkt, 1)
			byteChan <- packet{
				buffer[:size],
			}
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			log.Println("about to close the server")
			close(byteChan)
			l.Close()
		}
	}()

	return byteChan
}
