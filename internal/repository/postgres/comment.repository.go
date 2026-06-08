package postgres

import (
	"context"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type commentRepository struct {
    db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *commentRepository {
    return &commentRepository{db}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
    return nil
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
    return nil, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, comment *domain.Comment) error {
    return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, id int) error {
    return nil
}