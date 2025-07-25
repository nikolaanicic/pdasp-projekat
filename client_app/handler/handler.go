package handler

import (
	channelinterface "clientapp/channel_interface"
	"clientapp/data"
	"clientapp/dto"
	"clientapp/jwt"
	"clientapp/models"
	"clientapp/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users              map[string]models.UserInfo
	installedChainCode map[string]string
}

func New() *Handler {
	return &Handler{
		users:              data.GetInitialUsers(),
		installedChainCode: data.GetInitialChainCode(),
	}
}

func (h *Handler) logOutEveryoneExcept(userId string) {
	for k, v := range h.users {
		if k == userId {
			continue
		}

		for key := range v.ChannelInterfaces {
			delete(v.ChannelInterfaces, key)
		}
	}
}

func (h *Handler) logInUser(userInfo *models.UserInfo) error {
	h.logOutEveryoneExcept(userInfo.UserID)

	for chcodename := range h.installedChainCode {
		chi, err := channelinterface.New(chcodename, h.installedChainCode[chcodename], userInfo.UserID, userInfo.Organization)
		if err != nil {
			return err
		}

		userInfo.ChannelInterfaces[chcodename] = chi
	}

	return nil
}

func (h *Handler) Login(ctx *gin.Context) {
	var userLoginDto dto.UserLoginDto

	if err := ctx.ShouldBind(&userLoginDto); err != nil {
		ctx.JSON(http.StatusBadRequest, missingUserIDError)
		log.Println("[SERVER] [ERROR]", err)
		return
	}

	userInfo, exists := h.users[userLoginDto.UserID]

	if !exists {
		ctx.JSON(http.StatusNotFound, userIDNotFoundError)
	}

	token, err := jwt.GenerateJWT(userInfo.UserID, userInfo.Role)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToGenerateTokenError)
		return
	}

	if err := h.logInUser(&userInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed to connect to the chain"})
		return
	}

	log.Println(userInfo)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) InitLedger(ctx *gin.Context) {
	userIdParam, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, jwt.BadRequestNoAuthParamsError)
		return
	}

	user_id := userIdParam.(string)
	userInfo := h.users[user_id]

	channel := ctx.Param("channel")
	if channel == "" {
		ctx.JSON(http.StatusBadRequest, missingChannelError)
		return
	}

	chainCodeId := h.installedChainCode[channel]

	wallet, err := utils.CreateWallet(user_id, userInfo.Organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToCreatePopulateWalletError)
		return
	}

	gateway, err := utils.ConnectToGateway(wallet, userInfo.Organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToConnectGateway)
		return
	}

	defer gateway.Close()

	network, err := gateway.GetNetwork(channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToGetGatewayNetwork)
		return
	}

	contract := network.GetContract(chainCodeId)
	log.Println("[HANDLER] [SUBMIT TX] InitLedger")

	if _, err := contract.SubmitTransaction("InitLedger"); err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToSubmitTx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ledger successfully initialized"})
}

func (h *Handler) GetAllProducts(ctx *gin.Context) {

}
