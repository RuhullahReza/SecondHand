package middleware

import (
	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		payload, ok := c.Get("payload")
		if !ok {
			c.Abort()
			return
		}
	
		role := payload.(helper.Payload).Role

		if role != "ADMIN" {
			err := helper.NewAuthorization("only admin can access")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}