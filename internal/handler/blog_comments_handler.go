package handler

import "github.com/gin-gonic/gin"

// BlogCommentsHandler mirrors BlogCommentsController.java (placeholder).
type BlogCommentsHandler struct{}

func NewBlogCommentsHandler() *BlogCommentsHandler {
	return &BlogCommentsHandler{}
}

func (h *BlogCommentsHandler) RegisterRoutes(r *gin.Engine) {
	r.Group("/blog-comments")
}
