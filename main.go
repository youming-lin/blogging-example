package main

import (
	"blogging-example/internal/handlers"

	"github.com/gin-gonic/gin"
)

// Gin starter example from https://go.dev/doc/tutorial/web-service-gin

func main() {
	router := gin.Default()

	router.GET("/api/posts", handlers.GetPosts)
	router.POST("/api/posts", handlers.CreatePost)
	router.GET("/api/posts/:id", handlers.GetPostById)
	router.POST("/api/posts/:id/comments", handlers.AddCommentToPost)

	router.Run("localhost:8080")
}
