package repositories

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dsn string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("❌ erro ao parsear configuração do banco: %v", err)
	}

	// Optimize connection pool settings for performance
	config.MaxConns = 25                      // Maximum number of connections
	config.MinConns = 5                       // Minimum number of connections
	config.MaxConnLifetime = time.Hour        // Maximum connection lifetime
	config.MaxConnIdleTime = 30 * time.Minute // Maximum idle time
	config.HealthCheckPeriod = time.Minute    // Health check period

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("❌ erro ao conectar no banco: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ erro ao testar conexão com banco: %v", err)
	}

	log.Println("✅ Conexão com banco de dados estabelecida com sucesso")
	return pool
}
