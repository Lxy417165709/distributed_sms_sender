package model

import (
	"time"
)


const MqMessage = "message"

type Certification struct {
	RemainSendTimes int
	LastSendTime time.Time
}


type ResultOfSend struct {
	IsSuccess bool
	Msg string
}

type Message struct{
	SenderId string
	ReceiverId string
	Content string
	SmsSenderId int
	OrderNum int
}
