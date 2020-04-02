package pubsub

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestPubSub_single(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err := p.Subscribe("single", func(msg interface{}) {
		require.Equal(t, "blah-blah", msg)
		wg.Done()
	})
	require.NoError(t, err)

	err = p.Publish("single", "blah-blah")
	require.NoError(t, err)

	wg.Wait()
}

func TestPubSub_nonBlockPublish(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(11)

	_, err := p.Subscribe("non-bock-topic", func(msg interface{}) {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	})
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 11; i++ {
			err = p.Publish("non-bock-topic", "pew-pew")
			assert.NoError(t, err)
		}
		close(done)
	}()

	select {
	case <-time.After(10 * time.Second):
		t.Fatal("publish method must not be blocked")
	case <-done:
		// ok
	}
}

func TestPubSub_multipleSubjects(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(2)

	_, err := p.Subscribe("sub1", func(msg interface{}) {
		require.Equal(t, "blah-blah-1", msg)
		wg.Done()
	})
	require.NoError(t, err)

	_, err = p.Subscribe("sub2", func(msg interface{}) {
		require.Equal(t, "blah-blah-2", msg)
		wg.Done()
	})
	require.NoError(t, err)

	err = p.Publish("sub1", "blah-blah-1")
	require.NoError(t, err)

	err = p.Publish("sub2", "blah-blah-2")
	require.NoError(t, err)

	wg.Wait()
}

func TestPubSub_multipleSubscribers(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wgFirst := sync.WaitGroup{}
	wgFirst.Add(1)

	_, err := p.Subscribe("multiple", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		wgFirst.Done()
	})
	require.NoError(t, err)

	wgSecond := sync.WaitGroup{}
	wgSecond.Add(1)

	_, err = p.Subscribe("multiple", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		wgSecond.Done()
	})
	require.NoError(t, err)

	err = p.Publish("multiple", "pew-pew")
	require.NoError(t, err)

	wgFirst.Wait()
	wgSecond.Wait()
}

func TestPubSub_slowpoke(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	samples := 100

	wgSlow := sync.WaitGroup{}
	wgSlow.Add(samples)
	slowCtx, slowCancel := context.WithCancel(context.Background())
	defer func() {
		slowCancel()
		wgSlow.Wait()
	}()

	_, err := p.Subscribe("slowpoke", func(msg interface{}) {
		defer wgSlow.Done()

		select {
		case <-slowCtx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
		}
	})
	require.NoError(t, err)

	fastWg := sync.WaitGroup{}
	fastWg.Add(samples)

	_, err = p.Subscribe("slowpoke", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		fastWg.Done()
	})
	require.NoError(t, err)

	for i := 0; i < samples; i++ {
		err = p.Publish("slowpoke", "pew-pew")
		require.NoError(t, err)
	}

	done := make(chan struct{})
	go func() {
		fastWg.Wait()
		close(done)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("publish blocks on slowpoke?")
	case <-done:
		// ok
	}
}

func TestPubSub_unsubscribe(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	s, err := p.Subscribe("unsubscribe", func(msg interface{}) {
		t.Error("first subscriber must not be called")
	})
	require.NoError(t, err)

	s.Unsubscribe()

	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err = p.Subscribe("unsubscribe", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		wg.Done()
	})
	require.NoError(t, err)

	err = p.Publish("unsubscribe", "pew-pew")
	require.NoError(t, err)

	wg.Wait()
}

func TestPubSub_sequencePublishers(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(10)

	_, err := p.Subscribe("topic", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		wg.Done()
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		err := p.Publish("topic", "pew-pew")
		require.NoError(t, err)
	}

	wg.Wait()
}

func TestPubSub_concurrentPublishers(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(10)

	_, err := p.Subscribe("topic", func(msg interface{}) {
		require.Equal(t, "pew-pew", msg)
		wg.Done()
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		go func() {
			err := p.Publish("topic", "pew-pew")
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestPubSub_msgOrder(t *testing.T) {
	p := NewPubSub()
	defer checkedClose(t, p)

	wg := sync.WaitGroup{}
	wg.Add(15)

	c := uint64(0)
	_, err := p.Subscribe("topic", func(msg interface{}) {
		expected := atomic.AddUint64(&c, 1)
		require.Equal(t, expected, msg)
		wg.Done()
	})
	require.NoError(t, err)

	for i := uint64(1); i < 11; i++ {
		if i == 6 {
			c := uint64(5)
			_, subErr := p.Subscribe("topic", func(msg interface{}) {
				expected := atomic.AddUint64(&c, 1)
				require.Equal(t, expected, msg)
				wg.Done()
			})
			require.NoError(t, subErr)
		}

		err = p.Publish("topic", i)
		require.NoError(t, err)
	}

	wg.Wait()
}

func TestPubSub_failAfterClose(t *testing.T) {
	p := NewPubSub()
	err := p.Close(context.Background())
	require.NoError(t, err)

	_, err = p.Subscribe("topic", func(msg interface{}) {})
	require.Error(t, err)

	err = p.Publish("topic", "pew-pew")
	require.Error(t, err)
}

func TestPubSub_close(t *testing.T) {
	p := NewPubSub()

	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := p.Subscribe("unsubscribe", func(msg interface{}) {
		select {
		case <-ctx.Done():
			// fast exit
			return
		default:
			time.Sleep(2 * time.Second)
			wg.Done()
		}
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		err = p.Publish("unsubscribe", "pew-pew")
		require.NoError(t, err)
	}

	// do a lot of work
	time.Sleep(1 * time.Second)

	done := make(chan struct{})
	go func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer closeCancel()

		err := p.Close(closeCtx)
		if err != nil && closeCtx.Err() == nil {
			assert.NoError(t, err)
		}
		close(done)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("close must respect context timed out")
	case <-done:
		cancel()
		wg.Wait()
	}
}

func TestPubSub_closeWaitsMessageDelivery(t *testing.T) {
	p := NewPubSub()

	wg := sync.WaitGroup{}

	_, err := p.Subscribe("q", func(msg interface{}) {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	})
	require.NoError(t, err)

	for i := 0; i < 11; i++ {
		wg.Add(1)
		err = p.Publish("q", "pew-pew")
		require.NoError(t, err)
	}
	checkedClose(t, p)
	wg.Wait()
}

func checkedClose(t *testing.T, c interface {
	Close(ctx context.Context) error
}) {
	require.NoError(t, c.Close(context.Background()))
}
