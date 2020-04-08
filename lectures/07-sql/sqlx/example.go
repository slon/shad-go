rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE id = $1", 1)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var value struct {
        ID int          `db:"id"`
        Name string     `db:"name"`
    }
    if err := rows.StructScan(&value); err != nil {
        log.Fatal(err)
    }

    log.Println(value)
}

if err = rows.Err(); err != nil {
    log.Fatal(err)
}
