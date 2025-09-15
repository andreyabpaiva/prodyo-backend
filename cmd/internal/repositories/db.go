package repositories

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dsn string) *pgxpool.Pool {
    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        log.Fatalf("❌ erro ao conectar no banco: %v", err)
    }
    return pool
}
