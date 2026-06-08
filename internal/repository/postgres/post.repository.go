package postgres

import (
	"context"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
    db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
    return &postRepository{db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) error {
    return nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*domain.Post, error) {
    return nil, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
    return nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
    return nil
}