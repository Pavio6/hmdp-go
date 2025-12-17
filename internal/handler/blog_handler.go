package handler

import (
	"hmdp-backend/internal/dto/result"
	"hmdp-backend/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/model"
	"hmdp-backend/internal/service"
	"hmdp-backend/internal/utils"
)

type BlogHandler struct {
	blogService *service.BlogService
	userService *service.UserService
}

func NewBlogHandler(blogSvc *service.BlogService, userSvc *service.UserService) *BlogHandler {
	return &BlogHandler{blogService: blogSvc, userService: userSvc}
}

func (h *BlogHandler) SaveBlog(ctx *gin.Context) {
	var blog model.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, result.Fail("invalid payload"))
		return
	}
	loginUser, b := middleware.GetLoginUser(ctx)
	if !b {
		return
	}
	if loginUser == nil {
		ctx.JSON(http.StatusUnauthorized, result.Fail("未登录"))
		return
	}
	blog.UserID = loginUser.ID
	if err := h.blogService.Create(ctx.Request.Context(), &blog); err != nil {
		ctx.JSON(http.StatusInternalServerError, result.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, result.OkWithData(blog.ID))
}

func (h *BlogHandler) LikeBlog(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, result.Fail("invalid id"))
		return
	}
	if err := h.blogService.IncrementLike(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, result.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, result.Ok())
}

func (h *BlogHandler) QueryMyBlog(ctx *gin.Context) {
	loginUser, b := middleware.GetLoginUser(ctx)
	if !b {
		return
	}
	if loginUser == nil {
		ctx.JSON(http.StatusUnauthorized, result.Fail("未登录"))
		return
	}
	page := utils.ParsePage(ctx.Query("current"), 1)
	blogs, err := h.blogService.QueryByUser(ctx.Request.Context(), loginUser.ID, page, utils.MAX_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, result.OkWithData(blogs))
}

func (h *BlogHandler) QueryHotBlog(ctx *gin.Context) {
	page := utils.ParsePage(ctx.Query("current"), 1)
	blogs, err := h.blogService.QueryHot(ctx.Request.Context(), page, utils.MAX_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, result.Fail(err.Error()))
		return
	}
	for i := range blogs {
		user, err := h.userService.FindByID(ctx.Request.Context(), blogs[i].UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, result.Fail(err.Error()))
			return
		}
		if user != nil {
			blogs[i].Name = user.NickName
			blogs[i].Icon = user.Icon
		}
	}
	ctx.JSON(http.StatusOK, result.OkWithData(blogs))
}
