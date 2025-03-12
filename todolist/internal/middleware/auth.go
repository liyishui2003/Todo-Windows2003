package middleware

import (
	"net/http"
	"strings"
	"todolist/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 从请求头中获取 Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. 解析 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//回调函数，Parse内部会调用这个函数来获取JWTSecret的值
			return []byte(config.Cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			ctx.Abort()
			return
		}

		// 3. 将用户 ID 存入 Context
		//这里token.Claims是interface{}类型，无法对任何字段进行访问，需要转化成jwt.MapClaims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//存入Context中，每个请求都有独立的Context，由key(其中包含user_id)来区分
			ctx.Set("user_id", claims["user_id"])
		}

		ctx.Next()
	}
}
