tx, err := db.BeginTx(ctx, nil)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback()

if _, err = tx.ExecContext(ctx, ...); err != nil {
    log.Fatal(err)
}

if err = tx.Commit(); err != nil {
    log.Fatal(err)
}
