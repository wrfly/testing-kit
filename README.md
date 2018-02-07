# Testing-kit

some tools for testing

[![Build Status](https://travis-ci.org/wrfly/testing-kit.svg?branch=master)](https://travis-ci.org/wrfly/testing-kit)
## Stress-Testing

### udp-repeater

```txt
➜ ./udp-repeater -h
Usage of ./udp-repeater:
  -c int
        client num, default is `runtime.NumCPU()` (default 4)
  -l int
        local port to listen (default 1234)
  -r string
        remote addr to send packet (default "localhost:12345")
  -t int
        how many times to repeat (default 10000)

➜ ./udp-repeater
2018/01/09 20:35:46 listenning on :1234
2018/01/09 20:35:52 receive: 64/s       send: 635011/s
2018/01/09 20:35:53 receive: 89/s       send: 888880/s
2018/01/09 20:35:54 receive: 89/s       send: 894374/s
2018/01/09 20:35:55 receive: 89/s       send: 864181/s
2018/01/09 20:35:56 receive: 91/s       send: 935310/s
```

### tcp-repeater

```txt
➜ ./tcp-repeater -h
Usage of ./tcp-repeater:
  -c int
        client num, default is `runtime.NumCPU()` (default 4)
  -l int
        local port to listen (default 1234)
  -r string
        remote addr to send packet (default "localhost:12345")
  -t int
        how many times to repeat (default 10000)

➜ ./tcp-repeater
2018/01/09 22:32:11 listenning on :1234
2018/01/09 22:32:16 receive: 1/s        send: 10000/s
2018/01/09 22:32:18 receive: 1/s        send: 10000/s
2018/01/09 22:32:19 receive: 8/s        send: 80000/s
^C2018/01/09 22:32:21 about to close the server
2018/01/09 22:32:21 quit
➜
```

### creazy-udp

```txt
➜ ./creazy-udp -h
Usage of ./creazy-udp:
  -c runtime.NumCPU()
        client num, default is runtime.NumCPU() (default 4)
  -d string
        data to send, empty by default
  -n uint
        total number, set to positive for counter mode (ignore rate)
  -r int
        send rate, <num>/s, set to -1 for ultimate (default 10000)
  -t string
        target addr to send packet (default "localhost:12345")

➜ ./creazy-udp -d "hello"
2018/02/07 11:54:29 send 'hello' to [localhost:12345] with [4] clients at 10000/s
2018/02/07 11:54:30 Send: 10000/s       Total sent 10000
2018/02/07 11:54:31 Send: 10000/s       Total sent 20000
2018/02/07 11:54:32 Send: 10000/s       Total sent 30000
^C2018/02/07 11:54:32 quit
2018/02/07 11:54:32 total sent 40000
➜

➜ ./creazy-udp -d hello -n 100
2018/02/07 12:09:37 send 'hello' to [localhost:12345] with [4] clients fot [100] times
2018/02/07 12:09:37 Send 100 packets, done
2018/02/07 12:09:37 quit
2018/02/07 12:09:37 total sent 100
➜ 
```

## Statistic

### packet-counter

```txt
➜ ./packet-collector -h
Usage of ./packet-collector:
  -l int
        local port to listen (default 12345)

➜ ./packet-collector
2018/02/07 12:11:04 listenning on UDP :12345
2018/02/07 12:11:04 listenning on TCP :12345
2018/02/07 12:11:06 UDP: 100/s  Total: 100
^C2018/02/07 12:11:08 quit
2018/02/07 12:11:08 total received: UDP [100];TCP [0]
```

## Utils

### tokenbucket

`import "github.com/wrfly/testing-kit/util/tokenbucket"` to use this token bucket lib.

Performence(`example/main.go` with *i7-7600U*):

```txt
2018/01/15 23:12:45 5s test with 1000/s
2018/01/15 23:12:50 used: 5.000055967 s
2018/01/15 23:12:50 take: 4972
2018/01/15 23:12:50 drop: 114913190
2018/01/15 23:12:50 range [100000000] test
2018/01/15 23:12:51 used: 0.24327525 s
2018/01/15 23:12:51 take: 242
2018/01/15 23:12:51 drop: 99999758
```
