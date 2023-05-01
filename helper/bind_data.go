package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	ut "github.com/go-playground/universal-translator"
)

func BindData(c *gin.Context, ut ut.Translator, req interface{}) bool {

	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		err := NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}


	if err := c.ShouldBind(req); err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []string


			for _, err := range errs {

				invalidArgs = append(invalidArgs, err.Translate(ut))
				
			}

			err := NewBadRequest("Invalid request parameters. See invalidArgs")

			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		fallBack := NewInternal()

		c.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}

	return true
}