package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

var (
	sentNum  uint64
	errorNum uint64
)

type cli struct {
	C http.Client

	url     string
	host    string
	rate    int
	method  string
	payload io.Reader
}

func (c *cli) resolvHost() string {
	return "127.0.0.1"
}

func (c *cli) do(ctx context.Context) {
	req, err := http.NewRequest(c.method, c.url, c.payload)
	if err != nil {
		log.Fatalf("create request error: %s", err)
	}
	req.Host = c.resolvHost()

	if c.rate <= 0 {
		// unlimited mode
		go func() {
			for ctx.Err() == nil {
				go func() {
					if _, err := c.C.Do(req); err != nil {
						atomic.AddUint64(&errorNum, 1)
					}
					atomic.AddUint64(&sentNum, 1)
				}()
			}
		}()
	} else {
		go func() {
			for ctx.Err() == nil {
				// token bucket
				// if !tokenbucket.Take(1) {
				// 	continue
				// }
				go func() {
					if _, err := c.C.Do(req); err != nil {
						atomic.AddUint64(&errorNum, 1)
					}
					atomic.AddUint64(&sentNum, 1)
				}()
			}
		}()
	}

	go func() {
		printStatus(ctx)
	}()

	return
}

func printStatus(ctx context.Context) {
	go func() {
		var n uint64
		for ctx.Err() == nil {
			n = atomic.LoadUint64(&sentNum)
			time.Sleep(time.Second)
			if diff := atomic.LoadUint64(&sentNum) - n; diff != 0 {
				log.Printf("sent [%d], err [%d], RPS: %d/s",
					atomic.LoadUint64(&sentNum),
					atomic.LoadUint64(&errorNum), diff)
			}
		}
	}()
}

func handleCookie(f io.Reader, URL string) http.CookieJar {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatalf("handle cookie error: %s", err)
	}

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "",
		Value:  "",
		Path:   "/",
		Domain: u.Host,
	}
	cookies = append(cookies, cookie)
	jar.SetCookies(u, cookies)

	return jar
}

func main() {
	// flags
	rate := flag.Int("r", 1000, `sending rate, <num>/s`)
	timeout := flag.Int("timeout", 5, "timeout for the client")
	target := flag.String("u", "http://localhost", "target url")
	method := flag.String("m", "GET", "method: GET|POST")
	postFile := flag.String("pf", "", "post file path")
	cookieFile := flag.String("cf", "", "cookie file path")
	flag.Parse()

	// check and modify
	*method = strings.ToUpper(*method)

	// prepare
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	var (
		pld       io.Reader
		cookieJar http.CookieJar
	)
	if *postFile != "" {
		if f, err := os.Open(*postFile); err != nil {
			log.Fatalf("open file error: %s", err)
		} else {
			pld = f
		}
	}
	if *cookieFile != "" {
		if f, err := os.Open(*cookieFile); err != nil {
			log.Fatalf("open file error: %s", err)
		} else {
			cookieJar = handleCookie(f, *target)
		}
	}

	client := &cli{
		C: http.Client{
			Timeout: time.Duration(*timeout) * time.Second,
			Jar:     cookieJar,
		},
		url:     *target,
		rate:    *rate,
		method:  *method,
		payload: pld,
	}

	// ready
	log.Printf("send [%d]/s packages to [%s], method: %s, timeout: %d",
		*rate, *target, *method, *timeout)

	// go
	go client.do(ctx)

	<-sigChan

	log.Println("about to quit")
	cancel()

	runtime.GC()
	debug.FreeOSMemory()

	return
}
