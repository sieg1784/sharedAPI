package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("access_token")
		if token == "" {
			c.JSON(200, gin.H{
				"systemCode":    403,
				"systemMessage": "未攜帶token, 無存取權限",
				"data":          nil,
			})
			c.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParserToken(token)
		fmt.Println(claims)
		if err != nil {
			if err == TokenExpired {
				c.JSON(200, gin.H{
					"systemCode":    406,
					"systemMessage": "token已過期，請重新申請",
					"data":          nil,
				})
				c.Abort()
				return
			}

			c.JSON(400, gin.H{
				"systemCode":    400,
				"systemMessage": err.Error(),
				"data":          nil,
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)

	}
}
