package alerting

import (
	"time"
	"github.com/benbjohnson/clock"
)

type Ticker struct {
	C			chan time.Time
	clock		clock.Clock
	last		time.Time
	offset		time.Duration
	newOffset	chan time.Duration
}

func NewTicker(last time.Time, initialOffset time.Duration, c clock.Clock) *Ticker {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t := &Ticker{C: make(chan time.Time), clock: c, last: last, offset: initialOffset, newOffset: make(chan time.Duration)}
	go t.run()
	return t
}
func (t *Ticker) run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for {
		next := t.last.Add(time.Duration(1) * time.Second)
		diff := t.clock.Now().Add(-t.offset).Sub(next)
		if diff >= 0 {
			t.C <- next
			t.last = next
			continue
		}
		select {
		case <-t.clock.After(-diff):
		case offset := <-t.newOffset:
			t.offset = offset
		}
	}
}
