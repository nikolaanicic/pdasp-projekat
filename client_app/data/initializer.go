package data

import (
	channelinterface "clientapp/channel_interface"
	"clientapp/models"
)

func GetInitialUsers() map[string]*models.UserInfo {

	return map[string]*models.UserInfo{
		"jj1": {UserID: "jj1", Organization: "org1", Role: models.USER, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
		"it1": {UserID: "it1", Organization: "org2", Role: models.USER, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
		"ou1": {UserID: "ou1", Organization: "org3", Role: models.USER, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
		"s1":  {UserID: "s1", Organization: "org1", Role: models.ADMIN, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
		"s2":  {UserID: "s2", Organization: "org2", Role: models.ADMIN, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
		"s3":  {UserID: "s3", Organization: "org3", Role: models.ADMIN, ChannelInterfaces: make(map[string]*channelinterface.ChannelInterace, 0)},
	}

}

func GetInitialChainCode() map[string]string {
	return map[string]string{
		"tradechannel1": "traderchaincode1",
		"tradechannel2": "traderchaincode2",
	}
}
