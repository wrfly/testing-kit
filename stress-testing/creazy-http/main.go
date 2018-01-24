package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
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
	payload []byte
}

func (c *cli) resolvHost() string {
	return "106.75.85.238"
}

func (c *cli) do(ctx context.Context) {
	// req, err := http.NewRequest(c.method, c.url, strings.NewReader("sss"))
	req, err := http.NewRequest(c.method, c.url, bytes.NewReader(c.payload))
	if err != nil {
		log.Fatalf("create request error: %s", err)
	}

	// req.URL.Host = fmt.Sprintf("%s:%s",
	// 	c.resolvHost(), req.URL.Port())

	go func() {
		tk := time.NewTicker(time.Second)
		defer tk.Stop()
		for ctx.Err() == nil {
			lastSent := atomic.LoadUint64(&sentNum)
			<-tk.C
			diff := atomic.LoadUint64(&sentNum) - lastSent
			log.Printf("sent [%d], err [%d], RPS: %d/s",
				atomic.LoadUint64(&sentNum),
				atomic.LoadUint64(&errorNum), diff)
		}
	}()

	biu := func() {
		resp, err := c.C.Do(req)
		if err != nil {
			fmt.Println(err)
			atomic.AddUint64(&errorNum, 1)
			return
		}
		resp.Body.Close()
		atomic.AddUint64(&sentNum, 1)
	}

	if c.rate <= 0 {
		for ctx.Err() == nil {
			go biu()
		}
	} else {
		tk := time.NewTicker(time.Second)
		defer tk.Stop()
		for ctx.Err() == nil {
			for i := 0; i < c.rate; i++ {
				go biu()
			}
			<-tk.C
		}
	}
}

func handleCookie(data []byte, URL string) http.CookieJar {
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
		pld       []byte
		err       error
		cookieJar http.CookieJar
	)

	if *postFile != "" {
		pld, err = ioutil.ReadFile(*postFile)
		if err != nil {
			log.Fatalf("open file error: %s", err)
		}
	}

	if *cookieFile != "" {
		data, err := ioutil.ReadFile(*cookieFile)
		if err != nil {
			log.Fatalf("open file error: %s", err)
		}
		cookieJar = handleCookie(data, *target)
	}

	client := &cli{
		C: http.Client{
			Timeout: time.Second * 5,
			Jar:     cookieJar,
		},
		url:     *target,
		rate:    *rate,
		method:  *method,
		payload: pld,
	}

	// ready
	log.Printf("send [%d]/s packages to [%s], method: %s",
		*rate, *target, *method)

	// go
	go client.do(ctx)

	<-sigChan

	log.Println("about to quit")
	cancel()

	runtime.GC()
	debug.FreeOSMemory()

	return
}
