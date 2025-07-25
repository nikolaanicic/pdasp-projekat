package jwt

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_TOKEN_SECRET"))

func GenerateJWT(userId, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
	})

	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

func ExtractAndValidateToken(ctx *gin.Context) (*jwt.Token, error) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, unauthorizedTokenMissingError)
		return nil, jwt.ErrSignatureInvalid
	}

	token, err := ValidateJWT(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, unauthorizedTokenInvalidError)
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}
