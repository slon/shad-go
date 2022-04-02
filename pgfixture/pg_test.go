package pgfixture_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/pgfixture"
)

func TestLocalPostgres(t *testing.T) {
	dsn := pgfixture.Start(t)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	require.NoError(t, conn.Ping(ctx))

	_, err = conn.Exec(ctx, `CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT);`)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `INSERT INTO users (name) VALUES ($1);`, "Fedor")
	require.NoError(t, err)

	row := conn.QueryRow(ctx, `SELECT id, name FROM users LIMIT 1;`)

	var id int
	var name string

	require.NoError(t, row.Scan(&id, &name))
	require.Equal(t, 1, id)
	require.Equal(t, name, "Fedor")
}
