package types

import "zeroIM/pkg/constants"

// go get github.com/mitchellh/mapstructure

type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
)

type (
	Chat struct {
		ConversationId string             `mapstructure:"conversationId"`
		ChatType       constants.ChatType `json:"chatType"`
		SendId         string             `mapstructure:"sendId"`
		RecvId         string             `mapstructure:"recvId"`
		SendTime       int64              `mapstructure:"sendTime"`
		Msg            `mapstructure:"msg"`
	}
)
