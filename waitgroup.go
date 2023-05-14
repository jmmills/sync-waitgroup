// Package waitgroup implements a drop in replacement for the sync/waitgroup package
// that provides the ability to specify context.Context to use with Wait().
// By adding support for context.Context we can now specify a timeout or deadline for
// for a Wait.
package waitgroup

import (
	"sync/atomic"
)

type (
	// Interface defines an abstract interface for this waitgroup package.
	// Packages may take a dependency on this to allow for mocking and/or injection.
	Interface interface {
		Add(int32)
		Count() int32
		Done()
		Wait(...Option)
	}

	waitgroup struct {
		finished *atomic.Bool
		count    *atomic.Int32
		done     chan struct{}
	}
)

// New will initialize a waitgroup.
func New(delta ...int32) Interface {
	wg := &waitgroup{
		finished: new(atomic.Bool),
		count:    new(atomic.Int32),
		done:     make(chan struct{}),
	}

	for _, c := range delta {
		wg.Add(c)
	}

	return wg
}

// Add will increment the internal counter for waitgroup.
// Each go-routine that we are waiting for should call a Done for
// the sum of our waitgroup counter
func (wg *waitgroup) Add(n int32) {
	wg.count.Add(n)
}

// Count returns the current number of routines we are wainting for
func (wg *waitgroup) Count() int32 {
	return wg.count.Load()
}

// Done will mark a task as completed and decrement our counter.
// If Done is called more times than our wait counter (see Add),
// Done will panic as the waitgroup is expected to be finished.
func (wg *waitgroup) Done() {
	if wg.finished.Load() {
		return
	}

	wg.count.Add(-1)
	c := wg.count.Load()

	switch {
	case c > 0:
		return
	case c == 0:
		close(wg.done)
		wg.finished.Store(true)
		return
	}
}

// Wait will block until all tasks are completed.
// Optionally a WithContext option may be passed to allow
// for a timeout on this wait. If a timeout occurs this may
// be detected by checking the error value stored eithin the given
// context.
func (wg *waitgroup) Wait(opts ...Option) {
	opt := new(option)
	options(opts).apply(opt)

	if opt.WithContext != nil {
		for {
			select {
			case <-opt.WithContext.Done():
				return
			case <-wg.done:
				return
			}
		}
	} else {
		<-wg.done
	}

}

func (opts options) apply(opt *option) {
	for _, f := range opts {
		f(opt)
	}
}
