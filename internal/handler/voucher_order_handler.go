package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hmdp-backend/internal/dto"
)

// VoucherOrderHandler mirrors VoucherOrderController.java.
type VoucherOrderHandler struct{}

func NewVoucherOrderHandler() *VoucherOrderHandler {
	return &VoucherOrderHandler{}
}

func (h *VoucherOrderHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/voucher-order")
	group.POST("/seckill/:id", h.seckillVoucher)
}

func (h *VoucherOrderHandler) seckillVoucher(ctx *gin.Context) {
	if _, err := strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid voucher id"))
		return
	}
	ctx.JSON(http.StatusOK, dto.Fail("功能未完成"))
}
