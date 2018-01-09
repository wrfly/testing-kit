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
	"sync"
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
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()
		log.Printf("listenning on UDP %s\n", addr)
		serveUDP(ctx, *port)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		log.Printf("listenning on TCP %s\n", addr)
		serveTCP(ctx, *port)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		statistic(ctx)
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	cancel()

	wg.Wait()

	log.Println("quit")
}

func serveUDP(ctx context.Context, port int) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go func() {
		for {
			buffer := make([]byte, 1)
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
	var beforeUDP uint64
	var beforeTCP uint64
	for {
		beforeUDP = atomic.LoadUint64(&nUDP)
		beforeTCP = atomic.LoadUint64(&nTCP)
		time.Sleep(time.Second)
		if n := atomic.LoadUint64(&nUDP) - beforeUDP; n != 0 {
			log.Printf("UDP: %d/s\n", n)
		}
		if n := atomic.LoadUint64(&nTCP) - beforeTCP; n != 0 {
			log.Printf("TCP: %d/s\n", n)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
