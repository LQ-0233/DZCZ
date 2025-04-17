package middleware

import (
	"fabric-smart-evidence-storage/model"
	"fabric-smart-evidence-storage/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		username, role, err := util.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Set("role", role)

		// 继续处理请求
		c.Next()
	}
}

func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != nil {
			if role == model.AdminRole {
				c.Next()
				return
			}
		}
		c.JSON(200, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}

func InputAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != nil {
			if role == model.InputRole {
				c.Next()
				return
			}
		}
		c.JSON(200, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
