# Testing-kit

some tools for testing

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
  -n int
        send rate, <num>/s (default 10000)
  -r string
        remote addr to send packet (default "localhost:12345")

➜ ./creazy-udp -d "hello"
2018/01/11 22:27:02 send 'hello' to [localhost:12345] with [4] clients at 10000/s
2018/01/11 22:27:03 send: 10000/s
2018/01/11 22:27:04 send: 10000/s
2018/01/11 22:27:05 send: 10000/s
2018/01/11 22:27:06 send: 10000/s
2018/01/11 22:27:07 send: 10000/s
2018/01/11 22:27:08 send: 10000/s
2018/01/11 22:27:09 send: 10000/s
2018/01/11 22:27:10 send: 10000/s
2018/01/11 22:27:11 send: 10000/s
^C2018/01/11 22:27:11 quit
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
2018/01/09 21:28:26 listenning on UDP :12345
2018/01/09 21:28:26 listenning on TCP :12345
2018/01/09 21:28:33 TCP: 737038/s
2018/01/09 21:28:34 TCP: 2917838/s
2018/01/09 21:28:35 TCP: 3054689/s
2018/01/09 21:28:36 TCP: 3057794/s
2018/01/09 21:28:43 TCP: 2979672/s
2018/01/09 21:28:44 TCP: 2210824/s
^C2018/01/09 21:28:56 quit
```

## Utils

### tokenbucket

`import "github.com/wrfly/testing-kit/utils/tokenbucket"` to use this token bucket lib.

Performence(`examplt/main.go` with *i7-7600U*):

```txt
2018/01/13 02:16:43 5s test with 1000/s
2018/01/13 02:16:48 used: 5000064399 ns
2018/01/13 02:16:48 take: 4999
2018/01/13 02:16:48 drop: 57031252
2018/01/13 02:16:48 range [100000000] test
2018/01/13 02:16:53 used: 4942880532 ns
2018/01/13 02:16:53 take: 4942
2018/01/13 02:16:53 drop: 99995058
```