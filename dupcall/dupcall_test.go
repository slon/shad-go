package dupcall

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestCall_Simple(t *testing.T) {
	defer goleak.VerifyNone(t)

	called := 0

	var call Call
	result, err := call.Do(context.Background(), func(ctx context.Context) (interface{}, error) {
		called++
		return "ok", nil
	})

	require.NoError(t, err)
	require.Equal(t, result, "ok")
	require.Equal(t, called, 1)

	errFailed := errors.New("failed")

	result, err = call.Do(context.Background(), func(ctx context.Context) (interface{}, error) {
		called++
		return nil, errFailed
	})

	require.Equal(t, err, errFailed)
	require.Nil(t, result)
	require.Equal(t, called, 2)
}

func TestCall_Dedup(t *testing.T) {
	defer goleak.VerifyNone(t)

	called := 0
	cb := func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)

		called++
		return "ok", nil
	}

	var call Call
	for i := 0; i < 10; i++ {
		go call.Do(context.Background(), cb)
	}

	result, err := call.Do(context.Background(), cb)

	require.NoError(t, err)
	require.Equal(t, result, "ok")
	require.Equal(t, called, 1)
}

func TestCall_HalfCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	called := 0
	cb := func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)

		called++
		return "ok", nil
	}

	var call Call
	for i := 0; i < 10; i++ {
		go call.Do(context.Background(), cb)
	}

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		go call.Do(ctx, cb)

		time.Sleep(time.Millisecond)
		cancel()
	}

	result, err := call.Do(context.Background(), cb)

	require.NoError(t, err)
	require.Equal(t, result, "ok")
	require.Equal(t, called, 1)
}

func TestCall_FullCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	cancelled := make(chan struct{})
	cb := func(ctx context.Context) (interface{}, error) {
		<-ctx.Done()
		close(cancelled)

		return nil, nil
	}

	var call Call

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		go call.Do(ctx, cb)

		time.Sleep(time.Millisecond)
		cancel()
	}

	select {
	case <-cancelled:
		return

	case <-time.After(time.Millisecond * 100):
		t.Errorf("duplicate call not cancelled after 100ms")
	}
}

func TestCall_NonBlockingCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	var call Call
	cb := func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return nil, nil
	}

	cancelled := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_, err := call.Do(ctx, cb)
		assert.Error(t, err)
		close(cancelled)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-cancelled:
		return
	case <-time.After(50 * time.Millisecond):
		t.Errorf("cancelled call blocked for more that 50ms")
	}
}
