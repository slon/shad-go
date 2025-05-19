package dao

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.com/slon/shad-go/pgfixture"
)

func TestDao(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })
	dsn := pgfixture.Start(t)

	ctx := context.Background()

	dao, err := CreateDao(ctx, dsn)
	require.NoError(t, err)
	defer func() { _ = dao.Close() }()

	_, err = dao.Lookup(ctx, 42)
	require.ErrorIs(t, err, sql.ErrNoRows)

	aliceID, err := dao.Create(ctx, &User{Name: "Alice"})
	require.NoError(t, err)
	bobID, err := dao.Create(ctx, &User{Name: "Bob"})
	require.NoError(t, err)
	charlieID, err := dao.Create(ctx, &User{Name: "Charie"})
	require.NoError(t, err)

	require.Len(t, map[UserID]struct{}{aliceID: {}, bobID: {}, charlieID: {}}, 3)

	alice, err := dao.Lookup(ctx, aliceID)
	require.NoError(t, err)
	require.Equal(t, User{ID: aliceID, Name: "Alice"}, alice)

	require.NoError(t, dao.Delete(ctx, bobID))
	require.Error(t, dao.Delete(ctx, 999))

	_, err = dao.Lookup(ctx, bobID)
	require.ErrorIs(t, err, sql.ErrNoRows)

	require.NoError(t, dao.Update(ctx, &User{ID: charlieID, Name: "Chaplin"}))
	require.Error(t, dao.Update(ctx, &User{ID: 999, Name: "FooBar"}))

	users, err := dao.List(ctx)
	require.NoError(t, err)

	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	require.Equal(t, []User{
		{ID: aliceID, Name: "Alice"},
		{ID: charlieID, Name: "Chaplin"},
	}, users)
}
