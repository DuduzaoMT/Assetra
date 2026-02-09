package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

// connection struct to hold db connection
type Connection interface {
	Close() error
	DB() *pgx.Conn
}

type conn struct {
	database *pgx.Conn
}

func NewConnection(cfg Config) Connection {
	// fmt.Println("database url:", cfg.Dsn())

	ctx := context.Background()
	db, err := pgx.Connect(ctx, cfg.Dsn())
	if err != nil {
		panic(err)
	}
	return &conn{database: db}
}

func (c *conn) Close() error {
	return c.database.Close(context.Background())
}

func (c *conn) DB() *pgx.Conn {
	return c.database
}
