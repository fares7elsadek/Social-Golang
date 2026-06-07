package postgres

import (
   
    "github.com/fares7elsadek/Social-Golang/internal/domain"
)

type postRepository struct {
    
}

func NewPostRepository() *postRepository {
    return &postRepository{}
}

func (r *postRepository) CreatePost(post *domain.Post) error {
    return nil
}

func (r *postRepository) GetPostByID(id int) (*domain.Post, error) {
    return nil, nil
}

func (r *postRepository) UpdatePost(post *domain.Post) error {
    return nil
}

func (r *postRepository) DeletePost(id int) error {
    return nil
}