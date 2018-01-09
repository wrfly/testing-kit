# Testing-kit

some tools for testing

## strss-testing

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

## statistic

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