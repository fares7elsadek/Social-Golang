package service

import (
	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type commentService struct {
	commentRepo repository.CommentRepository
}


func NewCommentService(commentRepo repository.CommentRepository) *commentService {
	return &commentService{commentRepo: commentRepo}
}


func (s *commentService) CreateComment(authorID, postID int, content string) error {
	comment := &domain.Comment{
		AuthorID: authorID,
		PostID:   postID,
		Content:  content,
	}
	return s.commentRepo.CreateComment(comment)
}

func (s *commentService) GetCommentByID(id int) (*domain.Comment, error) {
	return s.commentRepo.GetCommentByID(id)
}

func (s *commentService) UpdateComment(id int, content string) error {
	comment := &domain.Comment{
		ID:      id,
		Content: content,
	}
	return s.commentRepo.UpdateComment(comment)
}

func (s *commentService) DeleteComment(id int) error {
	return s.commentRepo.DeleteComment(id)
}

