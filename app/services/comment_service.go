package services

import (
	"context"
	"errors"
	"strings"

	model "app/models"
	repository "app/repositories"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo, postRepo: postRepo}
}

func (s *CommentService) Create(ctx context.Context, postID int64, authorName, content string) (model.Comment, error) {
	if err := s.ensurePostExists(ctx, postID); err != nil {
		return model.Comment{}, err
	}

	return s.commentRepo.Create(ctx, model.Comment{
		PostID:     postID,
		AuthorName: strings.TrimSpace(authorName),
		Content:    strings.TrimSpace(content),
	})
}

func (s *CommentService) GetAllByPost(ctx context.Context, postID int64) ([]model.Comment, error) {
	if err := s.ensurePostExists(ctx, postID); err != nil {
		return nil, err
	}

	return s.commentRepo.GetAllByPostID(ctx, postID)
}

func (s *CommentService) ensurePostExists(ctx context.Context, postID int64) error {
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	return nil
}
