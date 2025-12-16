package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/model"
	"hmdp-backend/internal/service"
)

// VoucherHandler mirrors VoucherController.java.
type VoucherHandler struct {
	service service.VoucherService
}

func NewVoucherHandler(svc service.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: svc}
}

func (h *VoucherHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/voucher")
	group.POST("", h.addVoucher)
	group.POST("/seckill", h.addSeckillVoucher)
	group.GET("/list/:shopId", h.queryVoucherOfShop)
}

func (h *VoucherHandler) addVoucher(ctx *gin.Context) {
	var voucher model.Voucher
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid payload"))
		return
	}
	if err := h.service.Create(ctx.Request.Context(), &voucher); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(voucher.ID))
}

func (h *VoucherHandler) addSeckillVoucher(ctx *gin.Context) {
	var voucher model.Voucher
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid payload"))
		return
	}
	if err := h.service.AddSeckillVoucher(ctx.Request.Context(), &voucher); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(voucher.ID))
}

func (h *VoucherHandler) queryVoucherOfShop(ctx *gin.Context) {
	shopID, err := strconv.ParseInt(ctx.Param("shopId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid shop id"))
		return
	}
	vouchers, err := h.service.QueryVoucherOfShop(ctx.Request.Context(), shopID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(vouchers))
}
