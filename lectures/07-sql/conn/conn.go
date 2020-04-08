c, err := db.Conn(ctx)
if err != nil {
    log.Fatal(err)
}

_ := c.PingContext(ctx)
_, _ := c.QueryContext(...)
