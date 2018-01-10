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
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sentNum  uint64
	errorNum uint64
	jobWG    sync.WaitGroup
)

type cli struct {
	C http.Client

	url     string
	host    string
	times   int
	method  string
	payload io.Reader
}

func (c *cli) resolvHost() string {
	return "127.0.0.1"
}

func (c *cli) do(ctx context.Context) chan struct{} {
	req, err := http.NewRequest(c.method, c.url, c.payload)
	if err != nil {
		log.Fatalf("create request error: %s", err)
	}
	req.Host = c.resolvHost()

	done := make(chan struct{})

	if c.times <= 0 {
		go func() {
			for ctx.Err() == nil {
				jobWG.Add(1)
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
			for i := 0; i < c.times; i++ {
				if ctx.Err() != nil {
					break
				}
				jobWG.Add(1)
				go func() {
					defer jobWG.Done()
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
		// wait for all jobs done, if times > 0
		jobWG.Wait()
		close(done)
	}()

	return done
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
	times := flag.Int("t", 10000, `how many times do you need to send, 
        set to -1 for unlimited mode`)
	timeout := flag.Int("timeout", 5, "timeout for the client")
	target := flag.String("u", "http://localhost", "target url")
	method := flag.String("m", "GET", "method: GET|POST")
	postFile := flag.String("pf", "", "post file")
	cookieFile := flag.String("cf", "", "cookie file")
	flag.Parse()

	// check and modify
	*method = strings.ToUpper(*method)

	// prepare
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	// post file
	var pld io.Reader
	if *postFile != "" {
		if f, err := os.Open(*postFile); err != nil {
			log.Fatalf("open file error: %s", err)
		} else {
			pld = f
		}
	}
	var cookieJar http.CookieJar
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
		times:   *times,
		method:  *method,
		payload: pld,
	}

	// ready
	log.Printf("send [%d] packages to [%s], method: %s, timeout: %d",
		*times, *target, *method, *timeout)

	// go
	done := client.do(ctx)

	select {
	case <-sigChan:
		log.Println("canceled")
	case <-done:
		log.Println("all jobs done")
	}
	log.Println("about to quit")
	cancel()

	// quit with deadline
	tk := time.NewTicker(time.Second * 10)
	select {
	case <-tk.C:
		log.Println("force quit")
	case <-func() chan bool {
		c := make(chan bool, 1)
		jobWG.Wait()
		c <- true
		return c
	}():
		log.Println("quit")
	}

	return
}
