package resources

import (
	"context"
	"database/sql"
	"log"
)

func QueryDeadlock(ctx context.Context, db *sql.DB) {
	rows, _ := db.QueryContext(ctx, "SELECT id, name FROM users")
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		_ = rows.Scan(&id, &name)

		rowsAddrs, _ := db.QueryContext(
			ctx,
			"SELECT address FROM addresses WHERE user_id = $1",
			id,
		)
		defer rowsAddrs.Close()

		var addr string
		_ = rowsAddrs.Scan(&addr)

		log.Println(id, name, addr)
	}
}

func QueryDeadlockFixOne(ctx context.Context, db *sql.DB) {
	type Res struct {
		ID   int
		Name string
		Addr string
	}
	var values []Res
	rows, _ := db.QueryContext(ctx, "SELECT id, name FROM users")

	for rows.Next() {
		var res Res
		_ = rows.Scan(&res.ID, &res.Name)
		values = append(values, res)
	}
	rows.Close()

	for _, v := range values {
		_ = db.QueryRowContext(
			ctx, "SELECT address FROM addresses WHERE user_id = $1", v.ID,
		).Scan(&v.Addr)
		log.Println(v)
	}
}

func QueryDeadlockFixTwo(ctx context.Context, db *sql.DB) {
	rows, _ := db.QueryContext(
		ctx,
		"SELECT u.id, u.name, a.address FROM users AS u, addresses as a WHERE u.id == a.user_id",
	)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var addr string
		_ = rows.Scan(&id, &name, &addr)
		log.Println(id, name, addr)
	}
}
