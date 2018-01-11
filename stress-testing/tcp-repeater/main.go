package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	remote    *string
	times     *int
	port      *int
	clientNum *int

	connGroup    []net.Conn
	numPerClient int

	sentPkt    uint64 // send per second
	receivePkt uint64 // receive per second
)

type packet struct {
	data []byte
}

func newBuffer() interface{} {
	return make([]byte, 20480)
}

var bufPool = sync.Pool{
	New: newBuffer,
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
		conn, err := net.Dial("tcp", *remote)
		if err != nil {
			log.Printf("conn error: %s", err)
			return
		}
		connGroup = append(connGroup, conn)
	}

	go func() {
		numPerClient = *times / *clientNum
		pktChan := serveTCP(ctx, *port)
		for pkt := range pktChan {
			go sendTCP(ctx, pkt)
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

func sendTCP(ctx context.Context, pkt packet) {
	for _, conn := range connGroup {
		go func(conn net.Conn) {
			for i := 0; i < numPerClient; i++ {
				atomic.AddUint64(&sentPkt, 1)
				_, err := conn.Write(pkt.data)
				if err != nil {
					if err == io.EOF || ctx.Err() != nil {
						return
					}
					log.Fatal(err)
				}
			}
		}(conn)
	}
}

func serveTCP(ctx context.Context, port int) chan packet {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("listenning on %s\n", addr)

	pktChan := make(chan packet, 1000)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for ctx.Err() == nil {
			c, err := l.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Fatal(err)
			}
			go func(c net.Conn) {
				defer c.Close()
				for ctx.Err() == nil {
					reader := bufio.NewReader(c)
					for {
						bytes, err := reader.ReadBytes('\n') // ReadString('\n')
						if err != nil {
							c.Close()
							return
						}
						atomic.AddUint64(&receivePkt, 1)
						pktChan <- packet{
							bytes,
						}
					}
				}
			}(c)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			log.Println("about to close the server")
			close(pktChan)
			l.Close()
		}
	}()

	return pktChan
}
