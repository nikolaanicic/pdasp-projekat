package models

import (
	channelinterface "clientapp/channel_interface"
)

type UserInfo struct {
	UserID            string `json:"user_id"`
	Organization      string `json:"organization"`
	Role              string `json:"role"`
	ChannelInterfaces map[string]*channelinterface.ChannelInterace
}

func NewUserInfo(id string, org string, role string) UserInfo {
	return UserInfo{UserID: id, Organization: org, Role: role, ChannelInterfaces: map[string]*channelinterface.ChannelInterace{"tradechannel1": nil, "tradechannel2": nil}}
}
