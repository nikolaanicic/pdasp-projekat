package jwt

import "github.com/gin-gonic/gin"

var unauthorizedTokenMissingError = gin.H{"status": "unauthorized - token is missing"}
var unauthorizedTokenInvalidError = gin.H{"status": "unauthorized - token is invalid"}
var unauthorizedClaimsInvalidError = gin.H{"status": "unauthorized - claims are invalid"}
var unauthorizedRoleInvalidError = gin.H{"status": "unauthorized - no permission"}
var BadRequestNoAuthParamsError = gin.H{"status": "bad request - no auth params provided"}
