package handler

import "github.com/gin-gonic/gin"

var missingUserIDError = gin.H{"status": "bad-request - user id is required"}
var missingChannelError = gin.H{"status": "bad-request - channel is required"}

var userIDNotFoundError = gin.H{"status": "not found - user not found"}
var failedToGenerateTokenError = gin.H{"status": "internal server error - failed to generate token"}
var failedToCreatePopulateWalletError = gin.H{"status": "internal server error - failed to generate and populate wallet"}
var failedToConnectGateway = gin.H{"status": "internal server error - failed to connect gateway"}
var failedToGetGatewayNetwork = gin.H{"status": "internal server error - failed to get the gateway network"}
var failedToSubmitTx = gin.H{"status": "internal server error - failed to submit tx"}
