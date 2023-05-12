package waitgroup

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitGroup(t *testing.T) {

	t.Run("no context", func(t *testing.T) {
		wg := New(1)
		wg.Add(1)
		f := func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
		}

		go f()
		go f()
		wg.Wait()

		assert.Panics(t, wg.Done)
	})

	t.Run("WithContext no timeout", func(t *testing.T) {
		wg := New(2)
		f := func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
		}

		go f()
		go f()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg.Wait(WithContext(ctx))
		assert.NoError(t, wg.Err())
	})

	t.Run("WithContext timeout", func(t *testing.T) {
		tctx, stop := context.WithCancel(context.Background())
		wg := New(2)
		f := func() {
			defer wg.Done()
			select {
			case <-tctx.Done():
				t.Log("stopping")
				return
			}
		}

		go f()
		go f()
		actx, acancel := context.WithCancel(context.Background())
		acancel()
		wg.Wait(WithContext(actx))
		assert.Error(t, wg.Err())
		assert.ErrorIs(t, wg.Err(), context.Canceled)

		bctx, bcancel := context.WithTimeout(context.Background(), time.Second)
		defer bcancel()
		stop()
		wg.Wait(WithContext(bctx))
		assert.NoError(t, wg.Err())
	})
}
