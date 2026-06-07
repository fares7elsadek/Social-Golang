package postgres

import (
    
    "github.com/fares7elsadek/Social-Golang/internal/domain"
)

type commentRepository struct {
    
}

func NewCommentRepository() *commentRepository {
    return &commentRepository{}
}

func (r *commentRepository) CreateComment(comment *domain.Comment) error {
    return nil
}

func (r *commentRepository) GetCommentByID(id int) (*domain.Comment, error) {
    return nil, nil
}

func (r *commentRepository) UpdateComment(comment *domain.Comment) error {
    return nil
}

func (r *commentRepository) DeleteComment(id int) error {
    return nil
}