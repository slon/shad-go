package clickhouse

import (
	"context"
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go"
)

func Example(ctx context.Context) {
	db, _ := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	defer db.Close()

	// Начало батча
	tx, _ := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	// Описание батча
	stmt, _ := tx.PrepareContext(ctx, "INSERT INTO example (id) VALUES (?)")
	defer stmt.Close()

	// Добавление данных
	for i := 0; i < 100; i++ {
		_, _ = stmt.ExecContext(ctx, i)
	}

	// Отправка батча в ClickHouse
	_ = tx.Commit()
}
