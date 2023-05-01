package middleware

import (
	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("token")

		if err != nil {
			err := helper.NewAuthorization("token required")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		payload, ok := helper.ValidateToken(token)

		if !ok {
			err := helper.NewAuthorization("invalid token")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("payload", payload)

		c.Next()

	}
}