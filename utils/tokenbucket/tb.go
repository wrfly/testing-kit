/*
package tokenbucket provides a token bucket ...

example:

	bkt := New(1000, time.Second)
	for check() {
		if bkt.TakeOne() {
			take++
			continue
		}
		drop++
	}

*/
package tokenbucket

import (
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	lastGet int64
	period  int64
	enough  int32
}

var globalBucket *Bucket
var globalIniter = sync.Once{}

// New returns a token bucket with `tokens` per `period`
func New(tokens int64, period time.Duration) *Bucket {
	period /= time.Duration(tokens)
	return &Bucket{
		lastGet: time.Now().UnixNano(),
		period:  period.Nanoseconds(),
		enough:  1,
	}
}

// NewGlobal returns a token bucket which is
// global shared and thread-safe
func NewGlobal(tokens int64, period time.Duration) *Bucket {
	period /= time.Duration(tokens)
	globalIniter.Do(func() {
		globalBucket = &Bucket{
			lastGet: time.Now().UnixNano(),
			period:  period.Nanoseconds(),
			enough:  1,
		}
	})
	return globalBucket
}

// TakeOne return true if you can take 1 from the bucket
func (bkt *Bucket) TakeOne() bool {
	// same round: not enough
	if atomic.LoadInt32(&bkt.enough) == 0 {
		return false
	}
	atomic.StoreInt32(&bkt.enough, 0)

	// next round: enough
	diff := time.Now().UnixNano() - atomic.LoadInt64(&bkt.lastGet)
	n := diff / bkt.period

	go func() {
		time.Sleep(time.Duration(bkt.period - diff + n*bkt.period))
		atomic.StoreInt32(&bkt.enough, 1)
	}()

	if n == 0 {
		n = 1
	}
	atomic.AddInt64(&bkt.lastGet, bkt.period*n)

	return true
}
