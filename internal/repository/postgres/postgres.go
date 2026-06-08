package postgres


import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)


func Connect(connectionString string) (*pgxpool.Pool,error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil,err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return nil,err
	}

	return pool,nil
}