package postgres


import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(connectionString string) error {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return err
	}

	DB = pool
	return nil
}