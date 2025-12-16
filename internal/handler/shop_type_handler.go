package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/service"
)

// ShopTypeHandler mirrors ShopTypeController.java.
type ShopTypeHandler struct {
	service service.ShopTypeService
}

func NewShopTypeHandler(svc service.ShopTypeService) *ShopTypeHandler {
	return &ShopTypeHandler{service: svc}
}

func (h *ShopTypeHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/shop-type")
	group.GET("/list", h.queryTypeList)
}

func (h *ShopTypeHandler) queryTypeList(ctx *gin.Context) {
	types, err := h.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(types))
}
