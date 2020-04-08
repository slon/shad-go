name := sql.NullString{Value: "The Shrike",Valid: true}
db.ExecContext(
    ctx,
    "INSERT INTO users(name) VALUES(@name)"),
    sql.Named("name", name),
)
