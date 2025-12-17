package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/dto/result"
	"hmdp-backend/internal/utils"
)

const loginUserContextKey = "loginUser"

// LoginMiddleware 校验登录
func LoginMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 需要登录
		needAuth := !isAnonymousPath(ctx.Request.URL.Path)
		// 提取token
		token := extractToken(ctx)
		if token == "" {
			if needAuth {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, result.Fail("未登录"))
				return
			}
			ctx.Next()
			return
		}
		key := utils.LOGIN_USER_KEY + token
		// 从redis中获取用户信息
		data, err := rdb.HGetAll(ctx.Request.Context(), key).Result()
		if err != nil {
			if needAuth {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, result.Fail("登录验证失败"))
			} else {
				ctx.Next()
			}
			return
		}
		if len(data) == 0 {
			if needAuth {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, result.Fail("登录状态已失效"))
			} else {
				ctx.Next()
			}
			return
		}
		id, _ := strconv.ParseInt(data["id"], 10, 64)
		user := &dto.UserDTO{
			ID:       id,
			NickName: data["nickName"],
			Icon:     data["icon"],
		}
		ctx.Set(loginUserContextKey, user)
		// 刷新token有效期
		rdb.Expire(ctx, key, time.Duration(utils.LOGIN_USER_TTL)*time.Second)
		ctx.Next()
	}
}

// GetLoginUser 从 Gin Context 中读取登录用户信息
func GetLoginUser(ctx *gin.Context) (*dto.UserDTO, bool) {
	v, exists := ctx.Get(loginUserContextKey)
	if !exists {
		return nil, false
	}
	user, ok := v.(*dto.UserDTO)
	return user, ok
}

// isAnonymousPath 这些路径放行 不需要登录即可访问
func isAnonymousPath(path string) bool {
	for _, prefix := range []string{"/shop", "/voucher", "/shop-type", "/upload"} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	switch path {
	case "/blog/hot", "/user/code", "/user/login":
		return true
	default:
		return false
	}
}

// extractToken 提取token
func extractToken(ctx *gin.Context) string {
	token := ctx.GetHeader("authorization")
	if token == "" {
		token = ctx.Query("token")
	}
	return token
}
