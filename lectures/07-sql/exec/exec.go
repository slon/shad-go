res, err := db.ExecContext(ctx, "UPDATE users SET name = $1 WHERE id = $2", "William Mandella", 1)
if err != nil {
    log.Fatal(err)
}

lastID, _ := res.LastInsertId()
rowsAffeced, := res.RowsAffected()

log.Println(lastID, rowsAffected)
