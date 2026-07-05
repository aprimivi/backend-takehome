package services

import (
	"context"
	"errors"
	"strings"

	model "app/models"
	repository "app/repositories"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("you don't have permission to modify this post")
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(ctx context.Context, authorID int64, title, content string) (model.Post, error) {
	return s.postRepo.Create(ctx, model.Post{
		Title:    strings.TrimSpace(title),
		Content:  strings.TrimSpace(content),
		AuthorID: authorID,
	})
}

func (s *PostService) GetAll(ctx context.Context) ([]model.Post, error) {
	return s.postRepo.GetAll(ctx)
}

func (s *PostService) GetByID(ctx context.Context, id int64) (model.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return model.Post{}, ErrPostNotFound
		}
		return model.Post{}, err
	}

	return post, nil
}

func (s *PostService) Update(ctx context.Context, id, authorID int64, title, content string) (model.Post, error) {
	if err := s.checkOwnership(ctx, id, authorID); err != nil {
		return model.Post{}, err
	}

	updated, err := s.postRepo.Update(ctx, id, authorID, strings.TrimSpace(title), strings.TrimSpace(content))
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return model.Post{}, ErrPostNotFound
		}
		return model.Post{}, err
	}

	return updated, nil
}

func (s *PostService) Delete(ctx context.Context, id, authorID int64) error {
	if err := s.checkOwnership(ctx, id, authorID); err != nil {
		return err
	}

	if err := s.postRepo.Delete(ctx, id, authorID); err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	return nil
}

func (s *PostService) checkOwnership(ctx context.Context, id, authorID int64) error {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	if post.AuthorID != authorID {
		return ErrForbidden
	}

	return nil
}
