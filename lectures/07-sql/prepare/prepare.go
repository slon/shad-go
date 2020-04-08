stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

for i := 1; ; i++ {
    row, err := stmt.QueryRowContext(ctx, i)
    if err != nil {
        log.Fatal(err)
    }

    var name string
    if err = row.Scan(&name); err != nil {
        log.Fatal(err)
    }

    log.Println(name)
}
