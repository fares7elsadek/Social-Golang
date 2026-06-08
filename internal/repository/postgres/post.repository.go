package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepository struct {
    db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *postRepository {
    return &postRepository{db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) error {
    query := `INSERT INTO posts 
              (title, content, author_id) 
              VALUES ($1, $2, $3)
              RETURNING id`

    err := r.db.QueryRow(ctx,query,
            post.Title,
            post.Content,
            post.AuthorID,
        ).Scan(&post.ID)
    
    if err != nil {
		return fmt.Errorf("postRepository.CreatePost: %w", err)
	}

	return nil    
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*domain.Post, error) {
    post := &domain.Post{}
    query := `
            SELECT id,title,content,author_id,created_at,updated_at 
            FROM posts
            WHERE id= $1
        `
    err := r.db.QueryRow(ctx,query,
    id).Scan(&post.ID,&post.Title,&post.Content,&post.AuthorID,&post.CreatedAt,&post.UpdatedAt)
    
    if err != nil {
        if errors.Is(err,pgx.ErrNoRows){
            return nil, domain.ErrNotFound
        }
		return nil,fmt.Errorf("postRepository.GetPostByID: %w", err)
	}

    return post, nil
}


func (r *postRepository) GetPostsByAuthor(ctx context.Context, authorId, limit, offset int) ([]*domain.Post, error){
    query := `SELECT id, title, content, author_id, created_at, updated_at
              FROM posts
              WHERE author_id = $1
              ORDER BY created_at DESC
              LIMIT $2 OFFSET $3`
    
    rows, err := r.db.Query(ctx, query, authorId, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("postRepository.GetPostsByAuthor: %w", err)
    }

    defer rows.Close()

    posts := make([]*domain.Post, 0)
    for rows.Next(){
        p := &domain.Post{}
        if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
            return nil, fmt.Errorf("postRepository.GetPostsByAuthor: scan: %w", err)
        }
        posts = append(posts, p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("postRepository.GetPostsByAuthor: rows: %w", err)
    }
    return posts, nil
}



func (r *postRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
    query := `
            UPDATE posts
            SET title= $1, content= $2, updated_at = NOW() WHERE id= $3
        `
    result, err := r.db.Exec(ctx,query,
                    post.Title,
                    post.Content,
                    post.ID)
    

    if err != nil {
		return fmt.Errorf("postRepository.UpdatePost: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {

    query := `
            DELETE FROM posts
            WHERE id= $1
        `
    result, err := r.db.Exec(ctx,query,
                    id)
    

    if err != nil {
		return fmt.Errorf("postRepository.DeletePost: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}