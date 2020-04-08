db.ExecContext(
    ctx,
    "INSERT INTO users(name) VALUES(@name)"),
    sql.Named("name", "Rocinante"),
)
