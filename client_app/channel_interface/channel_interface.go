package channelinterface

import (
	"clientapp/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ChannelInterace struct {
	Wallet   *gateway.Wallet
	Gateway  *gateway.Gateway
	Contract *gateway.Contract
	Network  *gateway.Network
}

func New(channel string, chainCodeId string, userID string, organization string) (*ChannelInterace, error) {

	wallet, err := utils.CreateWallet(userID, organization)
	if err != nil {
		return nil, err
	}

	gateway, err := utils.ConnectToGateway(wallet, organization)
	if err != nil {
		return nil, err
	}

	network, err := gateway.GetNetwork(channel)
	if err != nil {
		return nil, err
	}

	contract := network.GetContract(chainCodeId)

	return &ChannelInterace{
		Wallet:   wallet,
		Gateway:  gateway,
		Contract: contract,
		Network:  network,
	}, nil
}

func (ci *ChannelInterace) Close() {
	ci.Gateway.Close()
	ci.Contract = nil
	ci.Wallet = nil
	ci.Network = nil
}
