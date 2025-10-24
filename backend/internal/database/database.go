package database

import (
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/jmoiron/sqlx"
)

// Connect opens a database connection using the provided DSN and configures sane pooling defaults.
func Connect(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxIdleTime(5 * time.Minute)
    db.SetConnMaxLifetime(60 * time.Minute)

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
