package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/repository"
)

type postService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository,userRepo repository.UserRepository) *postService {
	return &postService{postRepo: postRepo,userRepo: userRepo}
}


func (s *postService) CreatePost(ctx context.Context,userID int,title string ,content string) error {
	post := &domain.Post{
		AuthorID: userID,
		Title: title,
		Content:  content,
	}

	if _, err := s.userRepo.GetUserByID(ctx, post.AuthorID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("user not found: %w",domain.ErrNotFound)
		}
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	err := s.postRepo.CreatePost(ctx,post)
	if err!=nil {
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}

func (s *postService) GetPostByID(ctx context.Context,id int) (*domain.Post, error) {

	post, err := s.postRepo.GetPostByID(ctx,id)
	if err != nil {
		if errors.Is(err,domain.ErrNotFound){
			return nil,fmt.Errorf("post not found: %w",domain.ErrNotFound)
		}
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return post,nil

}

func (s *postService) GetPostsByAuthorID(ctx context.Context,authorId,limit,offset int) ([]*domain.Post, error){
	posts, err := s.postRepo.GetPostsByAuthor(ctx,authorId,limit,offset)
	
	if err != nil {
		return nil,fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return posts,nil
}

func (s *postService) UpdatePost(ctx context.Context,id int, postParams domain.UpdatePostParams) error {
	
	post,err := s.GetPostByID(ctx,id)
	
	if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return fmt.Errorf("post not found: %w", domain.ErrNotFound)
        }
        return fmt.Errorf("unexpected error: %w", domain.ErrServer)
    }

	if postParams.Title != nil {
		post.Title = *postParams.Title
	}

	if postParams.Content != nil {
		post.Content = *postParams.Content
	}

	if err := s.postRepo.UpdatePost(ctx,post); err!=nil{
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}

func (s *postService) DeletePost(ctx context.Context,id int) error {
	_,err := s.GetPostByID(ctx,id)

	if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return fmt.Errorf("post not found: %w", domain.ErrNotFound)
        }
        return fmt.Errorf("unexpected error: %w", domain.ErrServer)
    }

	if err := s.postRepo.DeletePost(ctx,id); err!=nil{
		return fmt.Errorf("unexpected error: %w", domain.ErrServer)
	}

	return nil
}	



