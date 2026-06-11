package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type commentService struct {
	commentRepo repository.CommentRepository
	userRepo repository.UserRepository
}


func NewCommentService(commentRepo repository.CommentRepository,userRepo repository.UserRepository) *commentService {
	return &commentService{commentRepo: commentRepo,userRepo: userRepo}
}


func (s *commentService) CreateComment(ctx context.Context,authorID, postID int, content string) error {
	comment := &domain.Comment{
		AuthorID: authorID,
		PostID:   postID,
		Content:  content,
	}
	

	if _, err := s.userRepo.GetUserByID(ctx, comment.AuthorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("user not found: %w",domain.ErrNotFound)
		}
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	err := s.commentRepo.CreateComment(ctx,comment)
	if err!=nil {
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}

func (s *commentService) GetCommentByID(ctx context.Context,id int) (*domain.Comment, error) {
	comment, err := s.commentRepo.GetCommentByID(ctx,id)
	if err != nil {
		if errors.Is(err,domain.ErrNotFound){
			return nil,fmt.Errorf("comment not found: %w",domain.ErrNotFound)
		}
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return comment,nil
}

func (s *commentService) GetCommentByPostID(ctx context.Context,postId,limit,offset int) ([]*domain.Comment, error){
	comments, err := s.commentRepo.GetCommentsByPostId(ctx,postId,limit,offset)
	
	if err != nil {
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return comments,nil
}

func (s *commentService) UpdateComment(ctx context.Context,id int, content string) error {
	comment,err := s.GetCommentByID(ctx,id)

	if err !=nil{
		if errors.Is(err,domain.ErrNotFound){
			return fmt.Errorf("comment not found: %w",domain.ErrNotFound)
		}
	}

	comment.Content = content

	if err:= s.commentRepo.UpdateComment(ctx,comment);err!=nil{
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}

func (s *commentService) DeleteComment(ctx context.Context,id int) error {
	_,err := s.GetCommentByID(ctx,id)
	if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return fmt.Errorf("comment not found: %w", domain.ErrNotFound)
        }
        return fmt.Errorf("unexpected error: %w", domain.ErrServer)
    }

	if err := s.commentRepo.DeleteComment(ctx,id); err!=nil{
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}

