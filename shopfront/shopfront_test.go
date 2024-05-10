package shopfront_test

import (
	"context"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/redisfixture"
	"gitlab.com/slon/shad-go/shopfront"
)

func TestShopfront(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()

	ctx := context.Background()

	c := shopfront.New(rdb)

	items, err := c.GetItems(ctx, []shopfront.ItemID{1, 2, 3, 4}, 42)
	require.NoError(t, err)
	require.Equal(t, []shopfront.Item{
		{},
		{},
		{},
		{},
	}, items)

	require.NoError(t, c.RecordView(ctx, 3, 42))
	require.NoError(t, c.RecordView(ctx, 2, 42))

	require.NoError(t, c.RecordView(ctx, 2, 4242))

	items, err = c.GetItems(ctx, []shopfront.ItemID{1, 2, 3, 4}, 42)
	require.NoError(t, err)
	require.Equal(t, []shopfront.Item{
		{},
		{ViewCount: 2, Viewed: true},
		{ViewCount: 1, Viewed: true},
		{},
	}, items)
}

func TestShopFrontConcurrent(t *testing.T) {
	goleak.VerifyNone(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(t),
	})
	defer func() { _ = rdb.Close() }()

	ctx := context.Background()
	c := shopfront.New(rdb)

	N := 10000
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			assert.NoError(t, c.RecordView(ctx, 1, 1))
			wg.Done()
		}()
	}
	wg.Wait()

	items, err := c.GetItems(ctx, []shopfront.ItemID{1}, 1)
	require.NoError(t, err)
	require.Equal(t, []shopfront.Item{
		{ViewCount: N, Viewed: true},
	}, items)
}

func BenchmarkShopfront(b *testing.B) {
	const nItems = 1024

	rdb := redis.NewClient(&redis.Options{
		Addr: redisfixture.StartRedis(b),
	})
	defer func() { _ = rdb.Close() }()

	ctx := context.Background()
	c := shopfront.New(rdb)

	var ids []shopfront.ItemID
	for i := 0; i < nItems; i++ {
		ids = append(ids, shopfront.ItemID(i))
		require.NoError(b, c.RecordView(ctx, shopfront.ItemID(i), 42))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetItems(ctx, ids, 42)
		require.NoError(b, err)
	}
}
