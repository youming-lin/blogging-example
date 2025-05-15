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
	CreatePost(ctx context.Context, blogPost *models.BlogPost) (*models.BlogPost, error)

	// Get all blog posts in database
	GetPosts(ctx context.Context) ([]models.BlogPost, error)

	// Get blog post by ID
	GetPostById(ctx context.Context, id int) (*models.BlogPost, error)
}

var cachedBlogPostRepository BlogPostRepository

type blogPostRespository struct {
	db *sql.DB
}

var _ BlogPostRepository = (*blogPostRespository)(nil)

// Get a blog post respository instance
func NewBlogPostRepository(uri string) (BlogPostRepository, error) {
	if cachedBlogPostRepository == nil {
		db, err := sql.Open("sqlite", uri)
		if err != nil {
			return nil, err
		}

		cachedBlogPostRepository = &blogPostRespository{
			db: db,
		}
	}

	return cachedBlogPostRepository, nil
}

// CreatePost implements BlogPostRepository.
func (b *blogPostRespository) CreatePost(ctx context.Context, blogPost *models.BlogPost) (*models.BlogPost, error) {
	statement := `
	INSERT INTO blog_post (title, content)
	VALUES (?, ?)
	RETURNING id
	`
	row := b.db.QueryRowContext(ctx, statement, blogPost.Title, blogPost.Content)

	if err := row.Scan(&blogPost.Id); err != nil {
		fmt.Printf("Scan: %v\n", err)
		return nil, err
		// return nil, fmt.Errorf("Cannot find post with ID %d", id)
	}

	return blogPost, nil
}

// GetPosts implements BlogPostRepository.
func (b *blogPostRespository) GetPosts(ctx context.Context) ([]models.BlogPost, error) {
	// Adapted from https://stackoverflow.com/a/17266044
	statement := "SELECT id, title, content FROM blog_post;"
	rows, err := b.db.QueryContext(ctx, statement)
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

// GetPostById implements BlogPostRepository.
func (b *blogPostRespository) GetPostById(ctx context.Context, id int) (*models.BlogPost, error) {
	statement := "SELECT id, title, content FROM blog_post WHERE id = ? LIMIT 1;"
	row := b.db.QueryRowContext(ctx, statement, id)
	ret := &models.BlogPost{}

	if err := row.Scan(&ret.Id, &ret.Title, &ret.Content); err != nil {
		fmt.Printf("Scan: %v\n", err)
		return nil, err
		// return nil, fmt.Errorf("Cannot find post with ID %d", id)
	}

	return ret, nil
}
