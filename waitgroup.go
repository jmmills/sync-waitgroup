package waitgroup

import (
	"sync"
	"sync/atomic"
)

type (
	Interface interface {
		Add(int32)
		Done()
		Wait(...Option)
		Err() error
	}

	waitgroup struct {
		m     sync.RWMutex // protects err
		err   error
		count *atomic.Int32
		done  chan struct{}
	}
)

func New(delta ...int32) Interface {
	wg := &waitgroup{
		count: new(atomic.Int32),
		done:  make(chan struct{}),
	}

	for _, c := range delta {
		wg.Add(c)
	}

	return wg
}

func (wg *waitgroup) Add(n int32) {
	wg.count.Add(n)
}

func (wg *waitgroup) Done() {
	wg.count.Add(-1)

	c := wg.count.Load()
	switch {
	case c > 0:
		return
	case c == 0:
		close(wg.done)
		return
	default:
		panic("Done() on a finished waitgroup")
	}
}

func (wg *waitgroup) Wait(opts ...Option) {
	opt := new(option)
	options(opts).apply(opt)

	if opt.WithContext != nil {
		wg.m.Lock()
		defer wg.m.Unlock()
		for {
			select {
			case <-opt.WithContext.Done():
				wg.err = opt.WithContext.Err()
				return
			case <-wg.done:
				wg.err = nil
				return
			}
		}
	} else {
		<-wg.done
	}

}

func (wg *waitgroup) Err() error {
	wg.m.RLock()
	defer wg.m.RUnlock()
	return wg.err
}

func (opts options) apply(opt *option) {
	for _, f := range opts {
		f(opt)
	}
}
