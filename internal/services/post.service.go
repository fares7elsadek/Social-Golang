package service

import (
	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *postService {
	return &postService{postRepo: postRepo}
}


func (s *postService) CreatePost(userID int, content string) error {
	post := &domain.Post{
		AuthorID: userID,
		Content:  content,
	}
	return s.postRepo.CreatePost(post)
}

func (s *postService) GetPostByID(id int) (*domain.Post, error) {
	return s.postRepo.GetPostByID(id)
}

func (s *postService) UpdatePost(id int, content string) error {
	post := &domain.Post{
		ID:      id,
		Content: content,
	}
	return s.postRepo.UpdatePost(post)
}

func (s *postService) DeletePost(id int) error {
	return s.postRepo.DeletePost(id)
}	



