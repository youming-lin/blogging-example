package handlers

import (
	"blogging-example/models"
	"blogging-example/repository"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var dburi = os.Getenv("DB_URI")

type ErrorResponse struct {
	Error string `json:"error"`
}

// Get all blog posts
func GetPosts(c *gin.Context) {
	repo, err := repository.NewBlogPostRepository(dburi)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	posts, err := repo.GetPosts(c)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

// Create a new blog post
func CreatePost(c *gin.Context) {
	newPost := &models.BlogPost{}
	if err := c.BindJSON(newPost); err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	repo, err := repository.NewBlogPostRepository(dburi)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	newPost, err = repo.CreatePost(c, newPost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPost)
}

// Get blog post by ID
func GetPostById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	repo, err := repository.NewBlogPostRepository(dburi)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	post, err := repo.GetPostById(c, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

// Add comment to post given post ID
func AddCommentToPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	newComment := &models.Comment{}
	if err := c.BindJSON(newComment); err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	newComment.PostId = id

	repo, err := repository.NewCommentRepository(dburi)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	newComment, err = repo.AddComment(c, newComment)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newComment)
}
