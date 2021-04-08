package conf

type DBStats struct {
    MaxOpenConnections int // Maximum number of open connections to the database; added in Go 1.11

    // Pool Status
    OpenConnections int // The number of established connections both in use and idle.
    InUse           int // The number of connections currently in use; added in Go 1.11
    Idle            int // The number of idle connections; added in Go 1.11

    // Counters
    WaitCount         int64         // The total number of connections waited for; added in Go 1.11
    WaitDuration      time.Duration // The total time blocked waiting for a new connection; added in Go 1.11
    MaxIdleClosed     int64         // The total number of connections closed due to SetMaxIdleConns; added in Go 1.11
    MaxIdleTimeClosed int64         // The total number of connections closed due to SetConnMaxIdleTime; added in Go 1.15
    MaxLifetimeClosed int64         // The total number of connections closed due to SetConnMaxLifetime; added in Go 1.11
}
