package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/util"
)

// SwaggerAuthMiddleware handles authentication and authorization for the swagger
func SwaggerDynamicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}

		// Generate the expected password dynamically (each request)
		expectedUsername := "AdmiN"
		expectedPassword := "admin" + util.GetFormattedDateTimeMinutesCombinedInStr(util.GetTimeNowByLoc())

		if user != expectedUsername || pass != expectedPassword {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}
