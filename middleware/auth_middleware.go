package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/util"
)

// AuthMiddleware handles authentication and authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Validate `customer_no` path variable
		customerNo := ctx.Param("customer_no")
		if customerNo == "" {
			util.WriteResponse(ctx, util.JSON, http.StatusBadRequest, []any{}, []any{}, "customer no is missing", 0, 0)
			ctx.Abort()
			return
		}

		// Extract token from Authorization header
		tokenString := util.GetTokenFromHeader(ctx)
		if tokenString == "" {
			util.WriteResponse(ctx, util.JSON, http.StatusUnauthorized, []any{}, []any{}, "Authorization header is missing or invalid", 0, 0)
			ctx.Abort()
			return
		}

		// Parse and validate token
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			util.WriteResponse(ctx, util.JSON, http.StatusUnauthorized, []any{}, []any{}, fmt.Sprintf("invalid token: %v", err), 0, 0)
			ctx.Abort()
			return
		}

		// Validate token expiration
		if util.IsTokenExpired(claims) {
			util.WriteResponse(ctx, util.JSON, http.StatusUnauthorized, []any{}, []any{}, "token has expired", 0, 0)
			ctx.Abort()
			return
		}

		// Validate role authorization
		// if !isAuthorized(ctx, claims["username"].(string)) {
		// 	util.WriteResponse(ctx, util.JSON, http.StatusForbidden, []any{}, "you have no access rights")
		// 	ctx.Abort()
		// 	return
		// }

		// Process the next middleware/handler
		ctx.Next()
	}
}
