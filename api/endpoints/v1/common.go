package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRedmine(secretToken string) gin.HandlerFunc {

	return func(c *gin.Context) {

		st, b := c.GetQueryArray("token")
		if b == true && len(st) > 0 {
			if st[0] == secretToken {
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
