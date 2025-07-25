package handler

import (
	channelinterface "clientapp/channel_interface"
	"clientapp/data"
	"clientapp/dto"
	"clientapp/jwt"
	"clientapp/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	users              map[string]*models.UserInfo
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

	if err := h.logInUser(userInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed to connect to the chain"})
		return
	}

	log.Println(userInfo.ChannelInterfaces["tradechannel1"])
	log.Println(userInfo.ChannelInterfaces["tradechannel2"])

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

	chi, ok := userInfo.ChannelInterfaces[channel]
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "channel doesn't exist"})
	}
	log.Println("[HANDLER] [SUBMIT TX] InitLedger")

	if _, err := chi.Contract.SubmitTransaction("InitLedger"); err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToSubmitTx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ledger successfully initialized"})
}

func (h *Handler) GetAllProducts(ctx *gin.Context) {
	userIdEntry, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, jwt.BadRequestNoAuthParamsError)
		return
	}

	user_id := userIdEntry.(string)
	userInfo := h.users[user_id]

	channel := ctx.Param("channel")
	if channel == "" {
		ctx.JSON(http.StatusBadRequest, missingChannelError)
		return
	}

	chi, ok := userInfo.ChannelInterfaces[channel]
	if !ok || chi == nil {
		log.Println(ok, chi)
		ctx.JSON(http.StatusNotFound, gin.H{"status": "channel doesn't exist"})
		return
	}

	log.Println("[HANDLER] [SUBMIT TX] GetAllProducts")
	response, err := chi.Contract.SubmitTransaction("GetAllProducts")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failedToSubmitTx)
		return
	}

	if response == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed to invoke the chain"})
	}

	var products []models.Product
	err = json.Unmarshal(response, &products)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) AddUser(ctx *gin.Context) {

	userIdEntry, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, jwt.BadRequestNoAuthParamsError)
		return
	}
	adminId := userIdEntry.(string)
	adminUserInfo := h.users[adminId]

	channel := ctx.Param("channel")
	if channel == "" {
		ctx.JSON(http.StatusBadRequest, missingChannelError)
		return
	}

	chi, ok := adminUserInfo.ChannelInterfaces[channel]
	if !ok || chi == nil {
		log.Println(ok, chi)
		ctx.JSON(http.StatusNotFound, gin.H{"status": "channel doesn't exist"})
		return
	}

	var user dto.UserCreateDto
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldn't resolve body"})
		return
	}

	if _, ok := h.users[user.ID]; ok {
		ctx.JSON(http.StatusConflict, gin.H{"status": "user already exists"})
		return
	}

	newUserInfo := models.UserInfo{
		UserID:            user.ID,
		Organization:      adminUserInfo.Organization,
		Role:              models.USER,
		ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0),
	}

	h.users[user.ID] = &newUserInfo
	log.Println("[HANDLER] [SUBMIT TX] CreateUser")

	newUser := models.User{
		ID:             user.ID,
		Name:           user.Name,
		LastName:       user.LastName,
		Email:          user.Email,
		AccountBalance: user.AccountBalance,
		ReceiptsID:     make([]string, 0),
	}

	newUserBytes, _ := json.Marshal(newUser)
	_, err := chi.Contract.SubmitTransaction("CreateUser", string(newUserBytes))

	if err != nil {
		log.Println("[ERROR]", err)
		ctx.JSON(http.StatusInternalServerError, failedToSubmitTx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "user created"})

}

func (h *Handler) BuyProduct(ctx *gin.Context) {
	userIdEntry, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, jwt.BadRequestNoAuthParamsError)
		return
	}
	user_id := userIdEntry.(string)
	userInfo := h.users[user_id]

	channel := ctx.Param("channel")
	if channel == "" {
		ctx.JSON(http.StatusBadRequest, missingChannelError)
		return
	}

	product_id := ctx.Param("product_id")
	if product_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "missing product_id"})
		return
	}

	chi, ok := userInfo.ChannelInterfaces[channel]
	if !ok || chi == nil {
		log.Println(ok, chi)
		ctx.JSON(http.StatusNotFound, gin.H{"status": "channel doesn't exist"})
		return
	}

	log.Println("[HANDLER] [SUBMIT TX] BuyProduct")
	_, err := chi.Contract.SubmitTransaction("BuyProduct", product_id, user_id)

	if err != nil {
		log.Println("[ERROR]", err)
		ctx.JSON(http.StatusInternalServerError, failedToSubmitTx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
