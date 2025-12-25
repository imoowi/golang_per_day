package middlewares

import (
	"fmt"
	"gindemo2/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				`error`: "Authorization header 不能为空",
			})
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != `bearer` {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				`error`: "Authorization header 的格式错误，必须是 Bearer {token}",
			})
			return
		}
		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				`error`:  "非法token",
				`detail`: err.Error(),
			})
			return
		}
		fmt.Println(`claims=`, claims)
		//设定当前登录的用户id，方便传下去
		c.Set(`user_id`, claims.UserID)
		c.Next()
	}
}
