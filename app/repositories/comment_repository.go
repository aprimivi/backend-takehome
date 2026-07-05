package repository

import (
	"context"
	"database/sql"
	"errors"

	model "app/models"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment model.Comment) (model.Comment, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO comments (post_id, author_name, content)
		VALUES (?, ?, ?)
	`, comment.PostID, comment.AuthorName, comment.Content)
	if err != nil {
		return model.Comment{}, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return model.Comment{}, err
	}

	return r.GetByID(ctx, insertedID)
}

func (r *CommentRepository) GetAllByPostID(ctx context.Context, postID int64) ([]model.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, post_id, author_name, content, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorName, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *CommentRepository) GetByID(ctx context.Context, id int64) (model.Comment, error) {
	var comment model.Comment
	err := r.db.QueryRowContext(ctx, `
		SELECT id, post_id, author_name, content, created_at
		FROM comments
		WHERE id = ?
	`, id).Scan(&comment.ID, &comment.PostID, &comment.AuthorName, &comment.Content, &comment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Comment{}, ErrCommentNotFound
		}

		return model.Comment{}, err
	}

	return comment, nil
}
