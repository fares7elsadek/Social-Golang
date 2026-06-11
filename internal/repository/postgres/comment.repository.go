package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type commentRepository struct {
    db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *commentRepository {
    return &commentRepository{db}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {

    query := `
            INSERT INTOR comments
            (content, author_id, post_id)
            VALUES($1, 2$, 3$)
            RETURNING id
        `
    
    err := r.db.QueryRow(ctx,query,
                comment.Content,
                comment.AuthorID,
                comment.PostID,).Scan(&comment.ID)
    
    if err!=nil{
        return fmt.Errorf("commentRepository.CreateComment: %w", err)
    }


    return nil
}

func (r *commentRepository) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
    query := `
            SELECT id, content, author_id, post_id, created_at, updated_at 
            FROM comments WHERE id = $1
        `
    comment := &domain.Comment{}

    err := r.db.QueryRow(ctx,query,
                id,).Scan(&comment.ID,&comment.Content,
                        &comment.AuthorID,&comment.PostID,&comment.CreatedAt,&comment.UpdatedAt)
    if err != nil {
        if errors.Is(err,pgx.ErrNoRows){
            return nil, domain.ErrNotFound
        }
		return nil,fmt.Errorf("commentRepository.GetCommentByID: %w", err)
	}       

    return nil, nil
}

func (r *commentRepository) GetCommentsByPostId(ctx context.Context, postId, limit, offset int) ([]*domain.Comment, error){
    query := `
            SELECT id, content, author_id, post_id, created_at, updated_at 
            FROM comments WHERE post_id = $1
            ORDER BY created_at DESC
            LIMIT $2 OFFSET $3
        `
    

    rows, err := r.db.Query(ctx,query,
                postId,limit,offset)

    if err != nil {
        return nil, fmt.Errorf("commentRepository.GetCommentsByAuthor: %w", err)
    }

    defer rows.Close()
    
    comments := make([]*domain.Comment,0)
    for rows.Next() {
        c := &domain.Comment{}
        if err := rows.Scan(&c.ID,&c.Content,&c.AuthorID,&c.PostID,&c.CreatedAt,&c.UpdatedAt); err != nil {
            return nil, fmt.Errorf("commentRepository.GetCommentsByAuthor: scan: %w", err)
        }
        comments = append(comments,c)
    }

     if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("commentRepository.GetCommentsByAuthor: rows: %w", err)
    }
    return comments, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, comment *domain.Comment) error {

    query := `
            UPDATE comments SET content= $1
            WHERE id = $1
        `
    result, err := r.db.Exec(ctx,query,comment.ID)

    if err != nil {
		return fmt.Errorf("commentRepository.UpdateComment: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, id int) error {
    query := `
            DELETE FROM comments WHERE id= $1
        `
    result, err := r.db.Exec(ctx,query,id)

    if err != nil {
		return fmt.Errorf("commentRepository.DeleteComment: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

    return nil
}