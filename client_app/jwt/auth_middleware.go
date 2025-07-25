package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ExtractAndValidateToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, unauthorizedTokenInvalidError)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, unauthorizedClaimsInvalidError)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"].(string))
		ctx.Set("role", claims["role"].(string))

		ctx.Next()
	}
}

func AuthorizationMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		providedRoleEntry, ok := ctx.Get("role")
		if !ok {
			ctx.JSON(http.StatusBadRequest, BadRequestNoAuthParamsError)
			return
		}
		providedRole := providedRoleEntry.(string)

		if providedRole != requiredRole {
			ctx.JSON(http.StatusUnauthorized, unauthorizedRoleInvalidError)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
