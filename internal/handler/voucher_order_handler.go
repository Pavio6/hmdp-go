package handler

import (
	"hmdp-backend/internal/dto/result"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoucherOrderHandler struct{}

func NewVoucherOrderHandler() *VoucherOrderHandler {
	return &VoucherOrderHandler{}
}

func (h *VoucherOrderHandler) SeckillVoucher(ctx *gin.Context) {
	if _, err := strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, result.Fail("invalid voucher id"))
		return
	}
	ctx.JSON(http.StatusOK, result.Fail("功能未完成"))
}
