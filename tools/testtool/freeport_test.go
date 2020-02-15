package testtool

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetFreePort(t *testing.T) {
	p, err := GetFreePort()
	require.NoError(t, err)
	require.NotEmpty(t, p)
}

func TestWaitForPort(t *testing.T) {
	p, err := GetFreePort()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	WaitForPort(ctx, p)

	require.Error(t, ctx.Err())
}
