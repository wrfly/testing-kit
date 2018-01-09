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
)

var (
	remote    *string
	times     *int
	port      *int
	wg        sync.WaitGroup
	connGroup []net.Conn
	clientNum *int
)

type packet struct {
	data []byte
}

func main() {
	port = flag.Int("l", 1111, "local port to listen")
	times = flag.Int("t", 100000, "how many times to repeat")
	remote = flag.String("r", "localhost:12345", "remote addr to send packet")
	clientNum = flag.Int("c", 10, "client num")
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
		wg.Add(1)
		defer wg.Done()
		pktChan := serveUDP(ctx, *port)
		for {
			select {
			case <-ctx.Done():
				for _, conn := range connGroup {
					conn.Close()
				}
				return
			case pkt := <-pktChan:
				go func(pkt packet) {
					for n, conn := range connGroup {
						log.Printf("[%d] send to %s\n", n, *remote)
						go func(conn net.Conn) {
							for i := 0; i < *times / *clientNum; i++ {
								conn.Write(pkt.data)
							}
						}(conn)
					}
				}(pkt)
			}
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	cancel()

	wg.Wait()

	log.Println("quit")

	// runtime.GC()
	// debug.FreeOSMemory()
	// fmem, _ := os.Create("leak.prof")
	// pprof.WriteHeapProfile(fmem)

}

func serveUDP(ctx context.Context, port int) chan packet {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("listenning on %s\n", addr)

	byteChan := make(chan packet, 1000)

	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer l.Close()
		buffer := make([]byte, 20480)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				size, source, err := l.ReadFrom(buffer)
				if err != nil {
					fmt.Print(err)
					continue
				}
				log.Printf("receive data[%d] from %s\n", size, source.String())
				byteChan <- packet{
					buffer[:size],
				}
			}
		}
	}()

	return byteChan
}
