/*
package tokenbucket provided a token bucket ...

example:



*/
package tokenbucket

import (
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	lastGet  int64
	period   int64
	numTook  int64
	numTotal int64
}

var globalBucket *Bucket
var globalIniter = sync.Once{}

// New returns a token bucket with `tokens` per `period`
func New(tokens int64, period time.Duration) *Bucket {
	period /= time.Duration(tokens)
	return &Bucket{
		lastGet:  time.Now().UnixNano(),
		period:   period.Nanoseconds(),
		numTook:  0,
		numTotal: 1,
	}
}

// NewGlobal returns a token bucket which is
// global shared and thread-safe
func NewGlobal(tokens int64, period time.Duration) *Bucket {
	globalIniter.Do(func() {
		globalBucket = &Bucket{
			lastGet:  time.Now().UnixNano(),
			period:   period.Nanoseconds(),
			numTook:  0,
			numTotal: tokens,
		}
	})
	return globalBucket
}

// TakeOne return true if you can take 1 from the bucket
func (bkt *Bucket) TakeOne() bool {
	return bkt.Take(1)
}

// Take return true if you can take `num` from the bucket
func (bkt *Bucket) Take(num ...int) bool {
	now := time.Now().UnixNano()
	if now > atomic.LoadInt64(&bkt.lastGet)+atomic.LoadInt64(&bkt.period) {
		atomic.StoreInt64(&bkt.numTook, 1)
		atomic.StoreInt64(&bkt.lastGet, now)
		return true
	} else {
		if atomic.AddInt64(&bkt.numTook, 1) <= bkt.numTotal {
			atomic.StoreInt64(&bkt.lastGet, now)
			return true
		}
	}
	return false
}
