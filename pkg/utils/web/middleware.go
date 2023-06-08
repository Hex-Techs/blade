package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/utils/token"
)

const CurrentUser = "user"

func GetCurrentUser(c *gin.Context) *token.Claims {
	return c.MustGet(CurrentUser).(*token.Claims)
}

// 登录验证中间件
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header中获取token
		t := c.Request.Header.Get("Authorization")
		if t == "" {
			c.JSON(http.StatusUnauthorized, ExceptResponse(http.StatusUnauthorized, "need login"))
			c.Abort()
			return
		}
		// 解析token
		claims, err := token.ParseJWTToken(t)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ExceptResponse(http.StatusUnauthorized, err))
			c.Abort()
			return
		}
		// 将用户信息保存到上下文中
		c.Set("user", claims)
		c.Next()
	}
}

// 必须是管理员才能访问的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetCurrentUser(c)
		if !claims.Admin {
			c.JSON(http.StatusForbidden, ExceptResponse(http.StatusForbidden, "no permission"))
			c.Abort()
			return
		}
		c.Next()
	}
}
