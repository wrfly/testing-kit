package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	addr string

	port *int
	nUDP uint64
	nTCP uint64
)

func main() {
	port = flag.Int("l", 12345, "local port to listen")
	flag.Parse()

	addr = fmt.Sprintf(":%d", *port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Printf("listenning on UDP %s\n", addr)
		serveUDP(ctx, *port)
	}()

	go func() {
		log.Printf("listenning on TCP %s\n", addr)
		serveTCP(ctx, *port)
	}()

	go func() {
		statistic(ctx)
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	cancel()

	log.Println("quit")
	log.Printf("total sent: UDP %v\tTCP%v\n",
		atomic.LoadUint64(&nUDP), atomic.LoadUint64(&nTCP))
}

func serveUDP(ctx context.Context, port int) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		buffer := make([]byte, 1)
		for {
			_, _, err := l.ReadFrom(buffer)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("error: %s\n", err)
				continue
			}
			atomic.AddUint64(&nUDP, 1)
		}
	}()

	<-ctx.Done()
	return
}

func serveTCP(ctx context.Context, port int) {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		buffer := make([]byte, 1)
		for {
			c, err := l.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Fatal(err)
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					_, err := c.Read(buffer)
					if err != nil {
						if err == io.EOF || ctx.Err() != nil {
							return
						}
						log.Printf("error: %s\n", err)
						continue
					}
					atomic.AddUint64(&nTCP, 1)
				}
			}(c)
		}
	}()

	<-ctx.Done()
	return
}

func statistic(ctx context.Context) {
	nu := atomic.LoadUint64(&nUDP)
	nt := atomic.LoadUint64(&nTCP)
	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			if n := atomic.LoadUint64(&nUDP) - nu; n != 0 {
				nu = atomic.LoadUint64(&nUDP)
				log.Printf("UDP: %v/s\tTotal: %v\n", n, nu)
			}
			if n := atomic.LoadUint64(&nTCP) - nt; n != 0 {
				nt = atomic.LoadUint64(&nTCP)
				log.Printf("TCP: %v/s\tTotal: %v\n", n, nt)
			}
		}
	}
}
