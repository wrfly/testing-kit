# Token bucket

`import "github.com/wrfly/testing-kit/utils/tokenbucket"` to use this token bucket lib.

This token bucket has no ticker to fill the bucket in background, just calculate the `next round` and make the bucket avaliable. Quite different from other token buckets implemented by locks or tickers or compare timestamps every time.

**But**, it has bugs in a single core host because of the `goroutine to reset the bucket`. You can use other implementions such as <https://github.com/tsenart/tb> or <https://github.com/bsm/ratelimit>

Performence(`examplt/main.go` with *i7-7600U*):

```txt
2018/01/15 23:12:45 5s test with 1000/s
2018/01/15 23:12:50 used: 5.000055967 s
2018/01/15 23:12:50 take: 4972
2018/01/15 23:12:50 drop: 114913190
2018/01/15 23:12:50 range [100000000] test
2018/01/15 23:12:51 used: 0.24327525 s
2018/01/15 23:12:51 take: 242
2018/01/15 23:12:51 drop: 99999758```
```

You can compare these token buckets at [compare](./compare).