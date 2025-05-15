package repository

import (
	"blogging-example/models"
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type CommentRepository interface {
	// Add a new comment to a blog post
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
}

var cachedCommentRepository CommentRepository

type commentRepository struct {
	db *sql.DB
}

var _ CommentRepository = (*commentRepository)(nil)

// Get a comment respository instance
func NewCommentRepository(uri string) (CommentRepository, error) {
	if cachedCommentRepository == nil {
		db, err := sql.Open("sqlite", uri)
		if err != nil {
			return nil, err
		}

		cachedCommentRepository = &commentRepository{
			db: db,
		}
	}

	return cachedCommentRepository, nil
}

// AddComment implements CommentRepository.
func (c *commentRepository) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	statement := `
	INSERT INTO comment (content, blog_post_id)
	VALUES (?, ?)
	RETURNING id
	`
	row := c.db.QueryRowContext(ctx, statement, comment.Content, comment.PostId)

	if err := row.Scan(&comment.Id); err != nil {
		fmt.Printf("Scan: %v\n", err)
		return nil, err
		// return nil, fmt.Errorf("Cannot find post with ID %d", id)
	}

	return comment, nil
}
