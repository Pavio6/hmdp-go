package handler

import "github.com/gin-gonic/gin"

// FollowHandler mirrors FollowController.java (placeholder).
type FollowHandler struct{}

func NewFollowHandler() *FollowHandler {
	return &FollowHandler{}
}

func (h *FollowHandler) RegisterRoutes(r *gin.Engine) {
	r.Group("/follow")
}
