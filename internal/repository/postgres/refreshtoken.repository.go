package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)


type refreshTokenRepo struct {
	db *pgxpool.Pool
}


func NewRefreshTokenRepo(db *pgxpool.Pool) *refreshTokenRepo {
	return &refreshTokenRepo{db:db}
}


func (t *refreshTokenRepo) Save(ctx context.Context, userID int, tokenID string, ttl time.Duration) error{

	query := `
			INSERT INTO refresh_token
			(tokenId,userId,ttl)
			VALUES($1, $2, $3)
		`
	
	_ , err := t.db.Exec(ctx,query,
				tokenID,
				userID,
				ttl,)
	
	if err != nil {
		return fmt.Errorf("Error setting refreshToken: %w",err)
	}

	return nil
}

func (t *refreshTokenRepo) Exists(ctx context.Context, userID int, tokenID string) (bool, error){

	query := `
        SELECT EXISTS(
            SELECT 1
            FROM refresh_token
            WHERE tokenId = $1
              AND userId = $2
        )
    `

    var exists bool

    err := t.db.QueryRow(
        ctx,
        query,
        tokenID,
        userID,
    ).Scan(&exists)

    if err != nil {
        return false, fmt.Errorf("checking refresh token existence: %w", err)
    }

    return exists, nil

}


func (t *refreshTokenRepo) Delete(ctx context.Context, userID int, tokenID string) error{
	query := `
			DELETE FROM refresh_token WHERE 
			tokenId = $1 AND userId = $2
	`

	result, err := t.db.Exec(ctx, query,
		tokenID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("refreshTokenRepo.Delete: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil

}

func (t *refreshTokenRepo) DeleteAll(ctx context.Context, userID int) error{
	query := `
			DELETE FROM refresh_token WHERE 
		    userId = $2
	`

	result, err := t.db.Exec(ctx, query,
		userID,
	)

	if err != nil {
		return fmt.Errorf("refreshTokenRepo.DeleteAll: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}


