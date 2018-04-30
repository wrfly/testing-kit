package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	addr string
	port *int

	nUDP, nTCP, sUDP, sTCP uint64
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
	log.Printf("total received: UDP [%d|%s];TCP [%d|%s]\n",
		atomic.LoadUint64(&nUDP), formatBytes(atomic.LoadUint64(&sUDP)),
		atomic.LoadUint64(&nTCP), formatBytes(atomic.LoadUint64(&sTCP)),
	)
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
			n, _, err := l.ReadFrom(buffer)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("error: %s\n", err)
				continue
			}
			atomic.AddUint64(&sUDP, uint64(n))
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
					n, err := c.Read(buffer)
					if err != nil {
						if err == io.EOF || ctx.Err() != nil {
							return
						}
						log.Printf("error: %s\n", err)
						continue
					}
					atomic.AddUint64(&sTCP, uint64(n))
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
	su := atomic.LoadUint64(&sUDP)
	st := atomic.LoadUint64(&sTCP)

	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			if n := atomic.LoadUint64(&nUDP) - nu; n != 0 {
				nu = atomic.LoadUint64(&nUDP)
				s := atomic.LoadUint64(&sUDP) - su
				size := formatBytes(s)
				log.Printf("UDP: %d/s|%s/s\tTotal: %d\n", n, size, nu)
			}
			if n := atomic.LoadUint64(&nTCP) - nt; n != 0 {
				nt = atomic.LoadUint64(&nTCP)
				s := atomic.LoadUint64(&sTCP) - st
				size := formatBytes(s)
				log.Printf("TCP: %d/s|%s/s\tTotal: %d\n", n, size, nt)
			}
		}
	}
}

// thanks to https://github.com/dustin/go-humanize/blob/master/bytes.go

func formatBytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}
