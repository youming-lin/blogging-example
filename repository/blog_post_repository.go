package repository

import (
	"blogging-example/models"
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Respository to CRUD blog posts
type BlogPostRepository interface {
	// Add a new blog post to database
	CreateBlogPost(ctx context.Context, blogPost models.BlogPost) error

	// Get all blog posts in database
	GetGlogPosts(ctx context.Context) ([]models.BlogPost, error)
}

var cache BlogPostRepository

type blogPostRespository struct {
	db *sql.DB
}

var _ BlogPostRepository = (*blogPostRespository)(nil)

// Get a blog post respository instance
func NewBlogPostRepository(uri string) (BlogPostRepository, error) {
	if cache == nil {
		db, err := sql.Open("sqlite", uri)
		if err != nil {
			return nil, err
		}

		cache = &blogPostRespository{
			db: db,
		}
	}

	return cache, nil
}

// CreateBlogPost implements BlogPostRepository.
func (b *blogPostRespository) CreateBlogPost(ctx context.Context, blogPost models.BlogPost) error {
	panic("unimplemented")
}

// GetGlogPosts implements BlogPostRepository.
func (b *blogPostRespository) GetGlogPosts(ctx context.Context) ([]models.BlogPost, error) {
	// Adapted from https://stackoverflow.com/a/17266044

	rows, err := b.db.QueryContext(ctx, "SELECT id, title, content FROM blog_post;")
	if err != nil {
		return nil, err
	}

	ret := []models.BlogPost{}

	for rows.Next() {
		var id int
		var title, content string

		if err = rows.Scan(&id, &title, &content); err != nil {
			fmt.Printf("Scan: %v\n", err)
		}

		ret = append(ret, models.BlogPost{
			Id:      id,
			Title:   title,
			Content: content,
		})
	}

	return ret, err
}
