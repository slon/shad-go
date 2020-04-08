for rows.Next() {
    var s sql.NullString
    if err := rows.Scan(&s); err != nil {
        log.Fatal(err)
    }

    if s.Valid {
       //
    } else {
       //
    }
}
