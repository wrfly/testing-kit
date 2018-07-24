/*
package tokenbucket provides a token bucket ...

example:

	bkt := NewSmoth(1000, time.Second)
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

type smothBucket struct {
	lastGet int64
	period  int64
	enough  int32
}

var globalBucket *smothBucket
var globalIniter = sync.Once{}

// NewSmoth returns a smoth token bucket with `tokens` per `period`
func NewSmoth(tokens int64, period time.Duration) *smothBucket {
	period /= time.Duration(tokens)
	return &smothBucket{
		lastGet: time.Now().UnixNano(),
		period:  period.Nanoseconds(),
		enough:  1,
	}
}

// NewGlobalSmoth returns a smoth token bucket which is
// global shared and thread-safe
func NewGlobalSmoth(tokens int64, period time.Duration) *smothBucket {
	period /= time.Duration(tokens)
	globalIniter.Do(func() {
		globalBucket = &smothBucket{
			lastGet: time.Now().UnixNano(),
			period:  period.Nanoseconds(),
			enough:  1,
		}
	})
	return globalBucket
}

// TakeOne return true if you can take 1 from the bucket
func (bkt *smothBucket) TakeOne() bool {
	// same round: not enough
	if atomic.LoadInt32(&bkt.enough) == 0 {
		return false
	}
	if atomic.SwapInt32(&bkt.enough, 0) == 0 {
		return false
	}

	// next round: enough
	diff := time.Now().UnixNano() - atomic.LoadInt64(&bkt.lastGet)
	n := diff / bkt.period

	if n == 0 {
		atomic.AddInt64(&bkt.lastGet, bkt.period)
	} else {
		atomic.AddInt64(&bkt.lastGet, bkt.period*n)
	}

	go func() {
		time.Sleep(time.Duration(bkt.period - diff + n*bkt.period))
		atomic.StoreInt32(&bkt.enough, 1)
	}()

	return true
}

type bucket struct {
	lastGet   int64
	total     int64
	toke      int64
	period    time.Duration
	first     time.Time
	inRound   int32
	notEnough int32
}

func New(tokens int64, period time.Duration) *bucket {
	return &bucket{
		lastGet: time.Now().UnixNano(),
		total:   tokens,
		period:  period,
	}
}

// TakeOne return true if you can take 1 from the bucket
func (bkt *bucket) TakeOne() bool {
	if atomic.LoadInt32(&bkt.notEnough) == 1 {
		// waiting for filling the bucket
		return false
	}

	if atomic.LoadInt32(&bkt.inRound) == 0 {
		bkt.first = time.Now()
		atomic.StoreInt32(&bkt.inRound, 1)
	}

	// enough
	if atomic.AddInt64(&bkt.toke, 1) <= bkt.total {
		return true
	}

	// not enough
	atomic.StoreInt32(&bkt.notEnough, 1)
	nextRoundTime := bkt.first.Add(bkt.period)
	sleep := nextRoundTime.Sub(time.Now())

	if sleep < 0 {
		atomic.StoreInt64(&bkt.toke, 0)
		atomic.StoreInt32(&bkt.inRound, 0)
		atomic.StoreInt32(&bkt.notEnough, 0)
	} else {
		go func() {
			time.Sleep(sleep)
			atomic.StoreInt64(&bkt.toke, 0)
			atomic.StoreInt32(&bkt.inRound, 0)
			atomic.StoreInt32(&bkt.notEnough, 0)
		}()
	}

	return false
}
