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

// ShopHandler mirrors ShopController.java behaviour.
type ShopHandler struct {
	service service.ShopService
}

func NewShopHandler(svc service.ShopService) *ShopHandler {
	return &ShopHandler{service: svc}
}

// RegisterRoutes binds shop endpoints.
func (h *ShopHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/shop")
	group.GET("/:id", h.queryShopByID)
	group.POST("", h.saveShop)
	group.PUT("", h.updateShop)
	group.GET("/of/type", h.queryShopByType)
	group.GET("/of/name", h.queryShopByName)
}

func (h *ShopHandler) queryShopByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid id"))
		return
	}
	shop, err := h.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(shop))
}

func (h *ShopHandler) saveShop(ctx *gin.Context) {
	var shop model.Shop
	if err := ctx.ShouldBindJSON(&shop); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid payload"))
		return
	}
	if err := h.service.Create(ctx.Request.Context(), &shop); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(shop.ID))
}

func (h *ShopHandler) updateShop(ctx *gin.Context) {
	var shop model.Shop
	if err := ctx.ShouldBindJSON(&shop); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid payload"))
		return
	}
	if err := h.service.Update(ctx.Request.Context(), &shop); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.Ok())
}

func (h *ShopHandler) queryShopByType(ctx *gin.Context) {
	typeIDStr := ctx.Query("typeId")
	if typeIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.Fail("typeId is required"))
		return
	}
	typeID, err := strconv.ParseInt(typeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid typeId"))
		return
	}
	page := utils.ParsePage(ctx.Query("current"), 1)
	shops, err := h.service.QueryByType(ctx.Request.Context(), typeID, page, utils.DEFAULT_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(shops))
}

func (h *ShopHandler) queryShopByName(ctx *gin.Context) {
	name := ctx.Query("name")
	page := utils.ParsePage(ctx.Query("current"), 1)
	shops, err := h.service.QueryByName(ctx.Request.Context(), name, page, utils.MAX_PAGE_SIZE)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(shops))
}
