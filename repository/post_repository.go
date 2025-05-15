package repository

import (
	"blogging-example/models"
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Respository to CRUD blog posts
type PostRepository interface {
	// Add a new blog post to database
	CreatePost(ctx context.Context, blogPost *models.BlogPost) (*models.BlogPost, error)

	// Get all blog posts in database
	GetPosts(ctx context.Context) ([]models.BlogPost, error)

	// Get blog post by ID
	GetPostById(ctx context.Context, id int) (*models.BlogPost, error)
}

// Cached object for reuse
var cachedBlogPostRepository PostRepository

type blogPostRespository struct {
	// DB connection
	db *sql.DB
}

var _ PostRepository = (*blogPostRespository)(nil)

// Get a blog post respository instance
func NewPostRepository(uri string) (PostRepository, error) {
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
	statement := `
	SELECT
		blog_post.id,
		blog_post.title,
		blog_post.content,
		COUNT(comment.id) AS number_of_comments
	FROM blog_post
	LEFT JOIN comment ON blog_post_id = blog_post.id
	GROUP BY 1, 2, 3
	;
	`
	rows, err := b.db.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}

	ret := []models.BlogPost{}

	// Adapted from https://stackoverflow.com/a/17266044
	for rows.Next() {
		var id int
		var title, content string
		var numberOfComments uint

		if err = rows.Scan(&id, &title, &content, &numberOfComments); err != nil {
			fmt.Printf("Scan: %v\n", err)
		}

		ret = append(ret, models.BlogPost{
			Id:               id,
			Title:            title,
			Content:          content,
			NumberOfComments: numberOfComments,
		})
	}

	return ret, err
}

// GetPostById implements BlogPostRepository.
func (b *blogPostRespository) GetPostById(ctx context.Context, id int) (*models.BlogPost, error) {
	statement := `
	SELECT 
		blog_post.id as id,
		blog_post.title as title,
		blog_post.content as content,
		comment.id as comment_id,
		comment.content as comment_content,
		comment.blog_post_id as comment_post_id
	FROM blog_post
	left join comment on blog_post_id = blog_post.id
	where blog_post.id = ?
	;
	`
	rows, err := b.db.QueryContext(ctx, statement, id)
	if err != nil {
		return nil, err
	}

	ret := &models.BlogPost{}

	for rows.Next() {
		comment := &models.Comment{}

		if err = rows.Scan(
			&ret.Id, &ret.Title, &ret.Content, // post data
			&comment.Id, &comment.Content, &comment.PostId, // comment data
		); err != nil {
			fmt.Printf("Scan: %v\n", err)
			return nil, err
		}

		ret.Comments = append(ret.Comments, comment)
	}

	return ret, nil
}
