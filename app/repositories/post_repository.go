package repository

import (
	"context"
	"database/sql"
	"errors"

	model "app/models"
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post model.Post) (model.Post, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO posts (title, content, author_id)
		VALUES (?, ?, ?)
	`, post.Title, post.Content, post.AuthorID)
	if err != nil {
		return model.Post{}, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return model.Post{}, err
	}

	return r.GetByID(ctx, insertedID)
}

func (r *PostRepository) GetAll(ctx context.Context) ([]model.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *PostRepository) GetByID(ctx context.Context, id int64) (model.Post, error) {
	var post model.Post
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		WHERE id = ?
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, ErrPostNotFound
		}

		return model.Post{}, err
	}

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, id, authorID int64, title, content string) (model.Post, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE posts
		SET title = ?, content = ?
		WHERE id = ? AND author_id = ?
	`, title, content, id, authorID)
	if err != nil {
		return model.Post{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return model.Post{}, err
	}
	if affected == 0 {
		return model.Post{}, ErrPostNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *PostRepository) Delete(ctx context.Context, id, authorID int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM posts
		WHERE id = ? AND author_id = ?
	`, id, authorID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrPostNotFound
	}

	return nil
}
