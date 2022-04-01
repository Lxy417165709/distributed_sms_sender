package model

import (
	"time"
)

const MqMessage = "message"

type Certification struct {
	RemainSendTimes int       // 剩余发送次数。
	LastSendTime    time.Time // 上次发送时间。
}

type ResultOfSend struct {
	IsSuccess bool
	Msg       string
}

type Message struct {
	SenderId    string
	ReceiverId  string
	Content     string
	SmsSenderId int
	OrderNum    int
}
