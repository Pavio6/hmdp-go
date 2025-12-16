package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/model"
	"hmdp-backend/internal/service"
	"hmdp-backend/internal/utils"
)

// BlogHandler mirrors BlogController.java.
type BlogHandler struct {
	blogService service.BlogService
	userService service.UserService
}

func NewBlogHandler(blogSvc service.BlogService, userSvc service.UserService) *BlogHandler {
	return &BlogHandler{blogService: blogSvc, userService: userSvc}
}

func (h *BlogHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/blog")
	group.POST("", h.saveBlog)
	group.PUT("/like/:id", h.likeBlog)
	group.GET("/of/me", h.queryMyBlog)
	group.GET("/hot", h.queryHotBlog)
}

func (h *BlogHandler) saveBlog(ctx *gin.Context) {
	var blog model.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid payload"))
		return
	}
	user := utils.GetUser()
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Fail("未登录"))
		return
	}
	blog.UserID = user.ID
	if err := h.blogService.Create(ctx.Request.Context(), &blog); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(blog.ID))
}

func (h *BlogHandler) likeBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid id"))
		return
	}
	if err := h.blogService.IncrementLike(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.Ok())
}

func (h *BlogHandler) queryMyBlog(ctx *gin.Context) {
	user := utils.GetUser()
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Fail("未登录"))
		return
	}
	page := utils.ParsePage(ctx.Query("current"), 1)
	blogs, err := h.blogService.QueryByUser(ctx.Request.Context(), user.ID, page, utils.MAX_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(blogs))
}

func (h *BlogHandler) queryHotBlog(ctx *gin.Context) {
	page := utils.ParsePage(ctx.Query("current"), 1)
	blogs, err := h.blogService.QueryHot(ctx.Request.Context(), page, utils.MAX_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	for i := range blogs {
		user, err := h.userService.FindByID(ctx.Request.Context(), blogs[i].UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
			return
		}
		if user != nil {
			blogs[i].Name = user.NickName
			blogs[i].Icon = user.Icon
		}
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(blogs))
}
