package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/service"
)

// UserHandler mirrors UserController.java.
type UserHandler struct {
	userService     service.UserService
	userInfoService service.UserInfoService
}

func NewUserHandler(userSvc service.UserService, userInfoSvc service.UserInfoService) *UserHandler {
	return &UserHandler{userService: userSvc, userInfoService: userInfoSvc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/user")
	group.POST("/code", h.sendCode)
	group.POST("/login", h.login)
	group.POST("/logout", h.logout)
	group.GET("/me", h.me)
	group.GET("/info/:id", h.info)
}

func (h *UserHandler) sendCode(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Fail("功能未完成"))
}

func (h *UserHandler) login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Fail("功能未完成"))
}

func (h *UserHandler) logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Fail("功能未完成"))
}

func (h *UserHandler) me(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Fail("功能未完成"))
}

func (h *UserHandler) info(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid id"))
		return
	}
	info, err := h.userInfoService.FindByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	if info == nil {
		ctx.JSON(http.StatusOK, dto.Ok())
		return
	}
	info.CreateTime = nil
	info.UpdateTime = nil
	ctx.JSON(http.StatusOK, dto.OkWithData(info))
}
