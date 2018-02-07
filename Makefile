.PHONY: statistic stress-testing util

packet-counter:
	go build -o bin/packet-counter statistic/packet-counter/main.go

statistic: packet-counter

creazy-http:
	go build -o bin/creazy-http stress-testing/creazy-http/main.go

creazy-udp:
	go build -o bin/creazy-udp stress-testing/creazy-udp/main.go

tcp-repeater:
	go build -o bin/tcp-repeater stress-testing/tcp-repeater/main.go

udp-repeater:
	go build -o bin/udp-repeater stress-testing/udp-repeater/main.go

stress-testing: creazy-http creazy-udp tcp-repeater udp-repeater

all: statistic stress-testing