package main

import (
	"blogging-example/internal/handlers"
	"blogging-example/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Gin starter example from https://go.dev/doc/tutorial/web-service-gin

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	repo, err := repository.NewBlogPostRepository("/Users/ylin/repos/blogging-example/data.sqlite")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	posts, err := repo.GetPosts(c)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.GET("/api/posts", handlers.GetPosts)
	router.POST("/api/posts", handlers.CreatePost)
	router.GET("/api/posts/:id", handlers.GetPostById)
	router.POST("/api/posts/:id/comments", handlers.AddCommentToPost)

	router.Run("localhost:8080")
}
