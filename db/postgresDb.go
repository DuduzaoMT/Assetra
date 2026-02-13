package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// connection struct to hold db connection pool
type Connection interface {
	Close()
	DB() *pgxpool.Pool
}

type conn struct {
	database *pgxpool.Pool
}

func NewConnection(cfg Config) Connection {
	// fmt.Println("database url:", cfg.Dsn())

	ctx := context.Background()
	// Create a connection pool instead of a single connection
	pool, err := pgxpool.Connect(ctx, cfg.Dsn())
	if err != nil {
		panic(err)
	}
	return &conn{database: pool}
}

func (c *conn) Close() {
	c.database.Close()
}

func (c *conn) DB() *pgxpool.Pool {
	return c.database
}
